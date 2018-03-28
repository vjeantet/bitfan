package parser

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/vjeantet/bitfan/core"
	"github.com/vjeantet/bitfan/entrypoint/parser/logstash"
)

var entryPointContent func(string, string, map[string]interface{}) ([]byte, string, error)

func parseConfigLocation(path string, options map[string]interface{}, pwd string, pickSections ...string) ([]core.Agent, error) {
	if path == "" {
		return []core.Agent{}, fmt.Errorf("no location provided to get content from ; options=%v ", options)
	}

	content, cwd, err := entryPointContent(path, pwd, options)

	if err != nil {
		return nil, err
	}

	agents, err := buildAgents(content, cwd, pickSections...)
	return agents, err
}

func BuildAgents(content []byte, pwd string, contentProvider func(string, string, map[string]interface{}) ([]byte, string, error)) ([]core.Agent, error) {
	entryPointContent = contentProvider
	return buildAgents(content, pwd)
}

func buildAgents(content []byte, pwd string, pickSections ...string) ([]core.Agent, error) {
	var i int
	agentConfList := []core.Agent{}
	if len(pickSections) == 0 {
		pickSections = []string{"input", "filter", "output"}
	}

	p := logstash.NewParser(bytes.NewReader(content))

	LSConfiguration, err := p.Parse()

	if err != nil {
		return agentConfList, err
	}

	outPorts := []core.Port{}

	if _, ok := LSConfiguration.Sections["input"]; ok && isInSlice("input", pickSections) {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["input"].Plugins); pluginIndex++ {
			plugin := LSConfiguration.Sections["input"].Plugins[pluginIndex]

			agents, tmpOutPorts, err := buildInputAgents(plugin, nil, pwd)
			if err != nil {
				return nil, err
			}

			agentConfList = append(agents, agentConfList...)
			outPorts = append(outPorts, tmpOutPorts...)
		}
	}

	if _, ok := LSConfiguration.Sections["filter"]; ok && isInSlice("filter", pickSections) {
		if _, ok := LSConfiguration.Sections["filter"]; ok {
			for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["filter"].Plugins); pluginIndex++ {
				var agents []core.Agent
				i++
				plugin := LSConfiguration.Sections["filter"].Plugins[pluginIndex]
				agents, outPorts, err = buildFilterAgents(plugin, outPorts, pwd)
				if err != nil {
					return nil, err
				}

				agentConfList = append(agents, agentConfList...)
			}
		}
	}

	if _, ok := LSConfiguration.Sections["output"]; ok && isInSlice("output", pickSections) {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["output"].Plugins); pluginIndex++ {
			var agents []core.Agent
			i++
			plugin := LSConfiguration.Sections["output"].Plugins[pluginIndex]
			agents, _, err = buildOutputAgents(plugin, outPorts, pwd)
			if err != nil {
				return nil, err
			}

			agentConfList = append(agents, agentConfList...)
		}
	}

	return agentConfList, nil
}

// TODO : this should return ports to be able to use multiple path use
func buildInputAgents(plugin *logstash.Plugin, lastOutPorts []core.Port, pwd string) ([]core.Agent, []core.Port, error) {
	agent := newAgent(plugin, pwd, "input_")

	// If agent is a "use"
	// build imported pipeline from path
	// connect import plugin Xsource to imported pipeline output
	if plugin.Name == "use" {
		if v, ok := agent.Options["path"]; ok {
			switch v.(type) {
			case string:
				agent.Options["path"] = []string{v.(string)}
				fileConfigAgents, err := parseConfigLocation(v.(string), agent.Options, pwd, "input", "filter")
				if err != nil {
					return nil, nil, err
				}

				// add agent "use" - set use agent Source as last From FileConfigAgents
				inPort := core.Port{AgentID: fileConfigAgents[0].ID, PortNumber: 0}
				agent.AgentSources = append(agent.AgentSources, inPort)
				fileConfigAgents = append([]core.Agent{agent}, fileConfigAgents...)

				outPort := core.Port{AgentID: fileConfigAgents[0].ID, PortNumber: 0}
				return fileConfigAgents, []core.Port{outPort}, nil
			case []interface{}:
				CombinedFileConfigAgents := []core.Agent{}
				newOutPorts := []core.Port{}
				for _, p := range v.([]interface{}) {
					// contruire le pipeline a
					fileConfigAgents, err := parseConfigLocation(p.(string), agent.Options, pwd, "input", "filter")
					if err != nil {
						return nil, nil, err
					}

					// save pipeline a for later return
					CombinedFileConfigAgents = append(CombinedFileConfigAgents, fileConfigAgents...)

					// add agent "use" - set use agent Source as last From FileConfigAgents
					inPort := core.Port{AgentID: fileConfigAgents[0].ID, PortNumber: 0}
					newOutPorts = append(newOutPorts, inPort)
				}

				// connect all collected inPort to "use" agent
				agent.AgentSources = append(agent.AgentSources, newOutPorts...)

				// add "use" plugin to combined pipelines
				CombinedFileConfigAgents = append([]core.Agent{agent}, CombinedFileConfigAgents...)

				// return  pipeline a b c ... with theirs respectives outputs
				return CombinedFileConfigAgents, []core.Port{{AgentID: agent.ID, PortNumber: 0}}, nil
			}
		}
	}

	// interval can be a number, a string number or a cron string pattern
	setAgentInterval(&agent)

	// @see commit dbeb4015a88893bffd6334d38f34f978312eff82
	setAgentTrace(&agent)

	setAgentPoolSize(&agent)

	outPort := core.Port{AgentID: agent.ID, PortNumber: 0}
	return []core.Agent{agent}, []core.Port{outPort}, nil
}

func buildOutputAgents(plugin *logstash.Plugin, lastOutPorts []core.Port, pwd string) ([]core.Agent, []core.Port, error) {
	agent := newAgent(plugin, pwd, "output_")

	// if its a use plugin
	// load filter and output parts of pipeline
	// connect pipeline Xsource to lastOutPorts
	// return pipelineagents with lastOutPorts intact
	// handle use plugin
	// If its a use agent
	// build the filter part of the pipeline
	// connect pipeline first agent Xsource to lastOutPorts output
	// return imported pipeline with its output
	if plugin.Name == "use" {
		if v, ok := agent.Options["path"]; ok {
			switch v.(type) {
			case string:
				agent.Options["path"] = []string{v.(string)}
				fileConfigAgents, err := parseConfigLocation(v.(string), agent.Options, pwd, "filter", "output")
				if err != nil {
					return nil, nil, err
				}

				firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
				for _, sourceport := range lastOutPorts {
					inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
					firstUsedAgent.AgentSources = append(firstUsedAgent.AgentSources, inPort)
				}

				//specific to output
				return fileConfigAgents, nil, nil

			case []interface{}:
				CombinedFileConfigAgents := []core.Agent{}
				for _, p := range v.([]interface{}) {
					fileConfigAgents, err := parseConfigLocation(p.(string), agent.Options, pwd, "filter", "output")
					if err != nil {
						return nil, nil, err
					}

					firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
					for _, sourceport := range lastOutPorts {
						inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
						firstUsedAgent.AgentSources = append(firstUsedAgent.AgentSources, inPort)
					}
					CombinedFileConfigAgents = append(CombinedFileConfigAgents, fileConfigAgents...)
				}
				// return  pipeline a b c ... with theirs respectives outputs
				return CombinedFileConfigAgents, nil, nil
			}
		}
	}

	// Plugin Sources
	agent.AgentSources = core.PortList{}
	for _, sourceport := range lastOutPorts {
		inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
		agent.AgentSources = append(agent.AgentSources, inPort)
	}

	agent_list := []core.Agent{}
	// Is this Plugin has conditional expressions ?
	if len(plugin.When) > 0 {
		var err error
		if agent_list, _, err = buildWhenBranch(&agent, plugin.When, "output"); err != nil {
			return nil, nil, err
		}
	}

	// interval can be a number, a string number or a cron string pattern
	setAgentInterval(&agent)

	// @see commit dbeb4015a88893bffd6334d38f34f978312eff82
	setAgentTrace(&agent)

	// ajoute l'agent à la liste des agents
	agent_list = append([]core.Agent{agent}, agent_list...)
	return agent_list, nil, nil
}

func buildFilterAgents(plugin *logstash.Plugin, lastOutPorts []core.Port, pwd string) ([]core.Agent, []core.Port, error) {
	agent := newAgent(plugin, pwd, "")
	agent.PoolSize = 2

	// handle use plugin
	// If its a use agent
	// build the filter part of the pipeline
	// connect pipeline first agent Xsource to lastOutPorts output
	// return imported pipeline with its output
	if plugin.Name == "use" {
		if v, ok := agent.Options["path"]; ok {
			switch v.(type) {
			case string:
				agent.Options["path"] = []string{v.(string)}
				fileConfigAgents, err := parseConfigLocation(v.(string), agent.Options, pwd, "filter")
				if err != nil {
					return nil, nil, err
				}

				firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
				for _, sourceport := range lastOutPorts {
					inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
					firstUsedAgent.AgentSources = append(firstUsedAgent.AgentSources, inPort)
				}

				newOutPorts := []core.Port{
					{AgentID: fileConfigAgents[0].ID, PortNumber: 0},
				}
				return fileConfigAgents, newOutPorts, nil

			case []interface{}:
				CombinedFileConfigAgents := []core.Agent{}
				newOutPorts := []core.Port{}
				for _, p := range v.([]interface{}) {
					// contruire le pipeline a
					fileConfigAgents, err := parseConfigLocation(p.(string), agent.Options, pwd, "filter")
					if err != nil {
						return nil, nil, err
					}

					// connect pipeline a first agent Xsource to lastOutPorts output
					firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
					for _, sourceport := range lastOutPorts {
						inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
						firstUsedAgent.AgentSources = append(firstUsedAgent.AgentSources, inPort)
					}
					// save pipeline a for later return
					CombinedFileConfigAgents = append(CombinedFileConfigAgents, fileConfigAgents...)
					// save pipeline a outputs for later return
					newOutPorts = append(newOutPorts, core.Port{AgentID: fileConfigAgents[0].ID, PortNumber: 0})
				}

				// connect all collected newOutPorts to "use" agent
				agent.AgentSources = append(agent.AgentSources, newOutPorts...)
				CombinedFileConfigAgents = append([]core.Agent{agent}, CombinedFileConfigAgents...)

				// return  pipeline a b c ... with theirs respectives outputs
				return CombinedFileConfigAgents, []core.Port{{AgentID: agent.ID, PortNumber: 0}}, nil
			}
		}
	}

	// route = set a pipeline, but do not reconnect it
	if plugin.Name == "route" {
		CombinedFileConfigAgents := []core.Agent{}
		for _, p := range agent.Options["path"].([]interface{}) {
			fileConfigAgents, err := parseConfigLocation(p.(string), agent.Options, pwd, "filter", "output")
			if err != nil {
				return nil, nil, err
			}

			// connect pipeline a last agent Xsource to lastOutPorts output
			lastUsedAgent := &fileConfigAgents[0]
			lastUsedAgent.AgentSources = append(lastUsedAgent.AgentSources, core.Port{AgentID: agent.ID, PortNumber: 0})

			CombinedFileConfigAgents = append(CombinedFileConfigAgents, fileConfigAgents...)
		}

		// connect route to lastOutPorts
		agent.AgentSources = append(agent.AgentSources, lastOutPorts...)
		// add route to routeedpipelines
		CombinedFileConfigAgents = append(CombinedFileConfigAgents, []core.Agent{agent}...)

		// return untouched outputsPorts
		return CombinedFileConfigAgents, []core.Port{{AgentID: agent.ID, PortNumber: 1}}, nil
	}

	// interval can be a number, a string number or a cron string pattern
	setAgentInterval(&agent)

	// @see commit dbeb4015a88893bffd6334d38f34f978312eff82
	setAgentTrace(&agent)

	setAgentPoolSize(&agent)

	// Plugin Sources
	agent.AgentSources = core.PortList{}
	for _, sourceport := range lastOutPorts {
		inPort := core.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
		agent.AgentSources = append(agent.AgentSources, inPort)
	}

	// By Default Agents output to port 0
	newOutPorts := []core.Port{
		{AgentID: agent.ID, PortNumber: 0},
	}

	agent_list := []core.Agent{}

	// Is this Plugin has conditional expressions ?
	if len(plugin.When) > 0 {
		var err error
		if agent_list, newOutPorts, err = buildWhenBranch(&agent, plugin.When, "filter"); err != nil {
			return nil, nil, err
		}
	}

	// ajoute l'agent à la liste des agents
	agent_list = append([]core.Agent{agent}, agent_list...)
	return agent_list, newOutPorts, nil
}

func isInSlice(needle string, candidates []string) bool {
	for _, symbolType := range candidates {
		if needle == symbolType {
			return true
		}
	}
	return false
}

func newAgent(plugin *logstash.Plugin, pwd string, labelPrefix string) core.Agent {
	var agent = core.Agent{}
	agent = core.NewAgent()
	agent.Type = labelPrefix + plugin.Name
	if plugin.Label == "" {
		agent.Label = plugin.Name
	} else {
		agent.Label = plugin.Label
	}

	agent.Buffer = 20
	agent.PoolSize = 1
	agent.Wd = pwd

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
	}

	// handle codecs
	if len(plugin.Codecs) > 0 {
		codecs := map[int]interface{}{}
		for i, codec := range plugin.Codecs {
			if codec.Name != "" {
				pcodec := core.NewCodec(codec.Name)
				for _, setting := range codec.Settings {
					pcodec.Options[setting.K] = setting.V
					if setting.K == "role" {
						pcodec.Role = setting.V.(string)
					}
				}

				codecs[i] = pcodec
			}
		}
		agent.Options["codecs"] = codecs
	}

	return agent
}

func setAgentInterval(agent *core.Agent) {
	interval := agent.Options["interval"]
	switch t := interval.(type) {
	case int, int8, int16, int32, int64:
		agent.Schedule = fmt.Sprintf("@every %ds", t)
	case string:
		if i, err := strconv.Atoi(t); err == nil {
			agent.Schedule = fmt.Sprintf("@every %ds", i)
		} else {
			agent.Schedule = t
		}
	}
}
func setAgentTrace(agent *core.Agent) {
	if trace, ok := agent.Options["trace"]; ok {
		switch t := trace.(type) {
		case string:
			agent.Trace = true
		case bool:
			agent.Trace = t
		}
	}
}
func setAgentPoolSize(agent *core.Agent) {
	if workers, ok := agent.Options["workers"]; ok {
		switch t := workers.(type) {
		case int64:
			agent.PoolSize = int(t)
		case int32:
			agent.PoolSize = int(t)
		case string:
			if i, err := strconv.Atoi(t); err == nil {
				agent.PoolSize = i
			}
		}
	}
	return
}

func buildWhenBranch(agent *core.Agent, Whens map[int]*logstash.When, sectionType string) ([]core.Agent, []core.Port, error) {
	agent_list := []core.Agent{}
	outPorts_when := []core.Port{}
	// le plugin WHEn est $plugin
	agent.Options["expressions"] = map[int]string{}
	elseOK := false
	// Loop over expressions in correct order
	for expressionIndex := 0; expressionIndex < len(Whens); expressionIndex++ {
		when := Whens[expressionIndex]
		//	enregistrer l'expression dans la conf agent
		agent.Options["expressions"].(map[int]string)[expressionIndex] = when.Expression
		if when.Expression == "true" {
			elseOK = true
		}
		// recupérer le outport associé (expressionIndex)
		expressionOutPorts := []core.Port{
			{AgentID: agent.ID, PortNumber: expressionIndex},
		}

		// construire les plugins associés à l'expression
		// en utilisant le outportA
		for pi := 0; pi < len(when.Plugins); pi++ {
			p := when.Plugins[pi]
			var agents []core.Agent
			var err error
			// récupérer le dernier outport du plugin créé il devient outportA
			if sectionType == "filter" {
				agents, expressionOutPorts, err = buildFilterAgents(p, expressionOutPorts, agent.Wd)
			} else if sectionType == "output" {
				agents, _, err = buildOutputAgents(p, expressionOutPorts, agent.Wd)
			}

			if err != nil {
				return nil, nil, err
			}

			// ajoute l'agent à la liste des agents
			agent_list = append(agents, agent_list...)
		}
		// ajouter le dernier outportA de l'expression au outport final du when
		outPorts_when = append(expressionOutPorts, outPorts_when...)
	}

	// If no else expression was found, insert one
	if elseOK == false {
		agent.Options["expressions"].(map[int]string)[len(agent.Options["expressions"].(map[int]string))] = "true"
		elseOutPorts := []core.Port{
			{AgentID: agent.ID, PortNumber: len(agent.Options["expressions"].(map[int]string)) - 1},
		}
		outPorts_when = append(elseOutPorts, outPorts_when...)
	}

	return agent_list, outPorts_when, nil
}
