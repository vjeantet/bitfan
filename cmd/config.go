package cmd

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/veino/config"
	"github.com/vjeantet/logstack/parser"
)

func parseConfig(logstackname string, content []byte) ([]config.Agent, error) {
	agentConfList := []config.Agent{}

	var i int

	p := parser.NewParser(bytes.NewReader(content))

	LSConfiguration, err := p.Parse()
	if err != nil {
		return agentConfList, err
	}

	outPorts := []config.Port{}

	for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["input"].Plugins); pluginIndex++ {
		i++
		plugin := LSConfiguration.Sections["input"].Plugins[pluginIndex]
		agent, outPort := buildInputAgent(logstackname, plugin, i)
		agentConfList = append([]config.Agent{agent}, agentConfList...)
		outPorts = append(outPorts, outPort)
	}

	if _, ok := LSConfiguration.Sections["filter"]; ok {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["filter"].Plugins); pluginIndex++ {
			var agents []config.Agent
			i++
			plugin := LSConfiguration.Sections["filter"].Plugins[pluginIndex]
			agents, outPorts, i = buildFilterAgents(logstackname, plugin, outPorts, i)
			agentConfList = append(agents, agentConfList...)
		}
	}

	for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["output"].Plugins); pluginIndex++ {
		var agents []config.Agent
		i++
		plugin := LSConfiguration.Sections["output"].Plugins[pluginIndex]
		agents, outPorts, i = buildOutputAgents(logstackname, plugin, outPorts, i)
		agentConfList = append(agents, agentConfList...)
	}

	return agentConfList, nil
}

func buildInputAgent(logstackname string, plugin *parser.Plugin, i int) (config.Agent, config.Port) {

	var agent config.Agent
	agent.Pipeline = logstackname
	agent.Type = "input_" + plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logstackname, plugin.Name, i)
	agent.Buffer = 200
	agent.PoolSize = 1

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
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

	return agent, config.Port{AgentName: agent.Name, PortNumber: 0}
}

func buildOutputAgents(logstackname string, plugin *parser.Plugin, lastOutPorts []config.Port, i int) ([]config.Agent, []config.Port, int) {
	agent_list := []config.Agent{}

	var agent config.Agent
	agent.Pipeline = logstackname
	agent.Type = "output_" + plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logstackname, plugin.Name, i)
	agent.Buffer = 200
	agent.PoolSize = 1

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
	}
	for _, codec := range plugin.Codecs {
		agent.Options["codec"] = codec.Name
	}

	// Plugin Sources
	agent.XSources = config.PortList{}
	for _, sourceport := range lastOutPorts {
		inPort := config.Port{AgentName: sourceport.AgentName, PortNumber: sourceport.PortNumber}
		agent.XSources = append(agent.XSources, inPort)
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
				config.Port{AgentName: agent.Name, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le expressionOutPorts
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent
				i++
				// récupérer le dernier outport du plugin créé il devient expressionOutPorts
				agents, _, i = buildOutputAgents(logstackname, p, expressionOutPorts, i)
				// ajoute l'agent à la liste des agents construits
				agent_list = append(agents, agent_list...)
			}
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list, lastOutPorts, i
}

func buildFilterAgents(logstackname string, plugin *parser.Plugin, lastOutPorts []config.Port, i int) ([]config.Agent, []config.Port, int) {

	agent_list := []config.Agent{}

	var agent config.Agent
	agent.Pipeline = logstackname
	agent.Type = plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logstackname, plugin.Name, i)
	agent.Buffer = 200
	agent.PoolSize = 2

	// Plugin configuration
	agent.Options = map[string]interface{}{}
	for _, setting := range plugin.Settings {
		agent.Options[setting.K] = setting.V
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
		inPort := config.Port{AgentName: sourceport.AgentName, PortNumber: sourceport.PortNumber}
		agent.XSources = append(agent.XSources, inPort)
	}

	// By Default Agents output to port 0
	newOutPorts := []config.Port{
		config.Port{AgentName: agent.Name, PortNumber: 0},
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
				config.Port{AgentName: agent.Name, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le outportA
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent
				i++
				// récupérer le dernier outport du plugin créé il devient outportA
				agents, expressionOutPorts, i = buildFilterAgents(logstackname, p, expressionOutPorts, i)
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
				config.Port{AgentName: agent.Name, PortNumber: len(agent.Options["expressions"].(map[int]string)) - 1},
			}
			newOutPorts = append(elseOutPorts, newOutPorts...)
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list, newOutPorts, i
}
