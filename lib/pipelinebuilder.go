package lib

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/veino/logfan/parser"
	"github.com/veino/veino/config"
)

func parseConfigLocation(name string, options map[string]interface{}, pwd string, pickSections ...string) ([]config.Agent, error) {
	var locs Locations

	if _, ok := options["path"]; ok {
		locs.Add(options["path"].(string), pwd)
	} else if _, ok := options["url"]; ok {
		locs.Add(options["url"].(string), pwd)
	} else {
		return []config.Agent{}, fmt.Errorf("no location provided to get content from ; options=%v ", options)
	}

	return locs.Items[0].configAgentsWithOptions(options, pickSections...)
}

func buildAgents(content []byte, pwd string, pickSections ...string) ([]config.Agent, error) {
	var i int
	agentConfList := []config.Agent{}
	if len(pickSections) == 0 {
		pickSections = []string{"input", "filter", "output"}
	}

	p := parser.NewParser(bytes.NewReader(content))

	LSConfiguration, err := p.Parse()

	if err != nil {
		return agentConfList, err
	}

	outPorts := []config.Port{}

	if _, ok := LSConfiguration.Sections["input"]; ok && isInSlice("input", pickSections) {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["input"].Plugins); pluginIndex++ {
			plugin := LSConfiguration.Sections["input"].Plugins[pluginIndex]

			agents := buildInputAgents(plugin, pwd)

			agentConfList = append(agents, agentConfList...)
			outPort := config.Port{AgentID: agents[0].ID, PortNumber: 0}
			outPorts = append(outPorts, outPort)
		}
	}

	if _, ok := LSConfiguration.Sections["filter"]; ok && isInSlice("filter", pickSections) {
		if _, ok := LSConfiguration.Sections["filter"]; ok {
			for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["filter"].Plugins); pluginIndex++ {
				var agents []config.Agent
				i++
				plugin := LSConfiguration.Sections["filter"].Plugins[pluginIndex]
				agents, outPorts = buildFilterAgents(plugin, outPorts, pwd)
				agentConfList = append(agents, agentConfList...)
			}
		}
	}

	if _, ok := LSConfiguration.Sections["output"]; ok && isInSlice("output", pickSections) {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["output"].Plugins); pluginIndex++ {
			var agents []config.Agent
			i++
			plugin := LSConfiguration.Sections["output"].Plugins[pluginIndex]
			agents = buildOutputAgents(plugin, outPorts, pwd)
			agentConfList = append(agents, agentConfList...)
		}
	}

	return agentConfList, nil
}

func buildInputAgents(plugin *parser.Plugin, pwd string) []config.Agent {

	var agent config.Agent
	agent = config.NewAgent()

	agent.Type = "input_" + plugin.Name
	agent.Label = fmt.Sprintf("%s", plugin.Name)
	agent.ID = agent.Label + "-" + agent.ID
	agent.Buffer = 200
	agent.PoolSize = 1
	agent.Wd = pwd

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
	}

	//todo : handle codec

	if plugin.Codec.Name != "" {
		agent.Options["codec"] = plugin.Codec.Name
	}

	// If agent is a "use"
	// build imported pipeline from path
	// connect import plugin Xsource to imported pipeline output
	if plugin.Name == "use" {
		fileConfigAgents, _ := parseConfigLocation("", agent.Options, pwd, "input", "filter")

		// add agent "use" - set use agent Source as last From FileConfigAgents
		inPort := config.Port{AgentID: fileConfigAgents[0].ID, PortNumber: 0}
		agent.XSources = append(agent.XSources, inPort)
		fileConfigAgents = append([]config.Agent{agent}, fileConfigAgents...)

		return fileConfigAgents
	}

	// interval can be a number, a string number or a cron string pattern
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

	return []config.Agent{agent}
}

func buildOutputAgents(plugin *parser.Plugin, lastOutPorts []config.Port, pwd string) []config.Agent {
	agent_list := []config.Agent{}

	var agent config.Agent
	agent = config.NewAgent()
	agent.Type = "output_" + plugin.Name
	agent.Label = fmt.Sprintf("%s", plugin.Name)
	agent.ID = agent.Label + "-" + agent.ID
	agent.Buffer = 200
	agent.PoolSize = 1
	agent.Wd = pwd

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
	}
	//todo : handle codec
	if plugin.Codec.Name != "" {
		agent.Options["codec"] = plugin.Codec.Name
	}
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
		fileConfigAgents, _ := parseConfigLocation("", agent.Options, pwd, "filter", "output")

		firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
		for _, sourceport := range lastOutPorts {
			inPort := config.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
			firstUsedAgent.XSources = append(firstUsedAgent.XSources, inPort)
		}

		//specific to output
		return fileConfigAgents
	}

	// Plugin Sources
	agent.XSources = config.PortList{}
	for _, sourceport := range lastOutPorts {
		inPort := config.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
		agent.XSources = append(agent.XSources, inPort)
	}

	if plugin.Codec != nil {
		agent.Options["codec"] = plugin.Codec.Name
	}

	// Is this Plugin has conditional expressions ?
	if len(plugin.When) > 0 {
		// outPorts_when := []port{}
		// le plugin WHEn est $plugin
		agent.Options["expressions"] = map[int]string{}
		// Loop over expressions in correct order
		for expressionIndex := 0; expressionIndex < len(plugin.When); expressionIndex++ {
			when := plugin.When[expressionIndex]
			//	enregistrer l'expression dans la conf agent
			agent.Options["expressions"].(map[int]string)[expressionIndex] = when.Expression

			// recupérer le outport associé (expressionIndex)
			expressionOutPorts := []config.Port{
				{AgentID: agent.ID, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le expressionOutPorts
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent

				// récupérer le dernier outport du plugin créé il devient expressionOutPorts
				agents = buildOutputAgents(p, expressionOutPorts, pwd)
				// ajoute l'agent à la liste des agents construits
				agent_list = append(agents, agent_list...)
			}
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list
}

func buildFilterAgents(plugin *parser.Plugin, lastOutPorts []config.Port, pwd string) ([]config.Agent, []config.Port) {

	agent_list := []config.Agent{}

	var agent config.Agent
	agent = config.NewAgent()

	agent.Type = plugin.Name
	agent.Label = fmt.Sprintf("%s", plugin.Name)
	agent.ID = agent.Label + "-" + agent.ID
	agent.Buffer = 200
	agent.PoolSize = 2
	agent.Wd = pwd

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
	}

	// handle use plugin
	// If its a use agent
	// build the filter part of the pipeline
	// connect pipeline first agent Xsource to lastOutPorts output
	// return imported pipeline with its output
	if plugin.Name == "use" {
		fileConfigAgents, _ := parseConfigLocation("", agent.Options, pwd, "filter")

		firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
		for _, sourceport := range lastOutPorts {
			inPort := config.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
			firstUsedAgent.XSources = append(firstUsedAgent.XSources, inPort)
		}

		newOutPorts := []config.Port{
			{AgentID: fileConfigAgents[0].ID, PortNumber: 0},
		}
		return fileConfigAgents, newOutPorts
	}

	// interval can be a number, a string number or a cron string pattern
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

	// Plugin Sources
	agent.XSources = config.PortList{}
	for _, sourceport := range lastOutPorts {
		inPort := config.Port{AgentID: sourceport.AgentID, PortNumber: sourceport.PortNumber}
		agent.XSources = append(agent.XSources, inPort)
	}

	// By Default Agents output to port 0
	newOutPorts := []config.Port{
		{AgentID: agent.ID, PortNumber: 0},
	}

	// Is this Plugin has conditional expressions ?
	if len(plugin.When) > 0 {
		outPorts_when := []config.Port{}
		// le plugin WHEn est $plugin
		agent.Options["expressions"] = map[int]string{}
		elseOK := false
		// Loop over expressions in correct order
		for expressionIndex := 0; expressionIndex < len(plugin.When); expressionIndex++ {
			when := plugin.When[expressionIndex]
			//	enregistrer l'expression dans la conf agent
			agent.Options["expressions"].(map[int]string)[expressionIndex] = when.Expression
			if when.Expression == "true" {
				elseOK = true
			}
			// recupérer le outport associé (expressionIndex)
			expressionOutPorts := []config.Port{
				{AgentID: agent.ID, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le outportA
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent
				// récupérer le dernier outport du plugin créé il devient outportA
				agents, expressionOutPorts = buildFilterAgents(p, expressionOutPorts, pwd)
				// ajoute l'agent à la liste des agents construits
				agent_list = append(agents, agent_list...)
			}
			// ajouter le dernier outportA de l'expression au outport final du when
			outPorts_when = append(expressionOutPorts, outPorts_when...)
		}
		newOutPorts = outPorts_when

		// If no else expression was found, insert one
		if elseOK == false {
			agent.Options["expressions"].(map[int]string)[len(agent.Options["expressions"].(map[int]string))] = "true"
			elseOutPorts := []config.Port{
				{AgentID: agent.ID, PortNumber: len(agent.Options["expressions"].(map[int]string)) - 1},
			}
			newOutPorts = append(elseOutPorts, newOutPorts...)
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list, newOutPorts
}

func isInSlice(needle string, candidates []string) bool {
	for _, symbolType := range candidates {
		if needle == symbolType {
			return true
		}
	}
	return false
}
