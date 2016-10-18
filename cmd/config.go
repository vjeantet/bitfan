package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/veino/logfan/parser"
	"github.com/veino/veino/config"
)

func parseConfig(logfanname string, content []byte, pwd string, pickSections ...string) ([]config.Agent, error) {
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
			i++
			plugin := LSConfiguration.Sections["input"].Plugins[pluginIndex]
			agent, outPort := buildInputAgent(logfanname, plugin, i, pwd)
			agentConfList = append([]config.Agent{agent}, agentConfList...)
			outPorts = append(outPorts, outPort)
		}
	}

	if _, ok := LSConfiguration.Sections["filter"]; ok && isInSlice("filter", pickSections) {
		if _, ok := LSConfiguration.Sections["filter"]; ok {
			for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["filter"].Plugins); pluginIndex++ {
				var agents []config.Agent
				i++
				plugin := LSConfiguration.Sections["filter"].Plugins[pluginIndex]
				agents, outPorts, i = buildFilterAgents(logfanname, plugin, outPorts, i, pwd)
				agentConfList = append(agents, agentConfList...)
			}
		}
	}

	if _, ok := LSConfiguration.Sections["output"]; ok && isInSlice("output", pickSections) {
		for pluginIndex := 0; pluginIndex < len(LSConfiguration.Sections["output"].Plugins); pluginIndex++ {
			var agents []config.Agent
			i++
			plugin := LSConfiguration.Sections["output"].Plugins[pluginIndex]
			agents, outPorts, i = buildOutputAgents(logfanname, plugin, outPorts, i, pwd)
			agentConfList = append(agents, agentConfList...)
		}
	}

	return agentConfList, nil
}

func buildInputAgent(logfanname string, plugin *parser.Plugin, i int, pwd string) (config.Agent, config.Port) {

	var agent config.Agent
	agent.Pipeline = logfanname
	agent.Type = "input_" + plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logfanname, plugin.Name, i)
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

	// If agent is a "use"
	// build imported pipeline from path
	// connect import plugin Xsource to imported pipeline output

	return agent, config.Port{AgentName: agent.Name, PortNumber: 0}
}

func buildOutputAgents(logfanname string, plugin *parser.Plugin, lastOutPorts []config.Port, i int, pwd string) ([]config.Agent, []config.Port, int) {
	agent_list := []config.Agent{}

	var agent config.Agent
	agent.Pipeline = logfanname
	agent.Type = "output_" + plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logfanname, plugin.Name, i)
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

	// if its a use plugin
	// load filter and output pars of pipeline
	// connect pipeline Xsource to lastOutPorts
	// return pipelineagents with lastOutPorts intact

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
				{AgentName: agent.Name, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le expressionOutPorts
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent
				i++
				// récupérer le dernier outport du plugin créé il devient expressionOutPorts
				agents, _, i = buildOutputAgents(logfanname, p, expressionOutPorts, i, pwd)
				// ajoute l'agent à la liste des agents construits
				agent_list = append(agents, agent_list...)
			}
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list, lastOutPorts, i
}

func buildFilterAgents(logfanname string, plugin *parser.Plugin, lastOutPorts []config.Port, i int, pwd string) ([]config.Agent, []config.Port, int) {

	agent_list := []config.Agent{}

	var agent config.Agent
	agent.Pipeline = logfanname
	agent.Type = plugin.Name
	agent.Name = fmt.Sprintf("%s_%s-%03d", logfanname, plugin.Name, i)
	agent.Buffer = 200
	agent.PoolSize = 2

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
		var ucwd = ""
		var content []byte

		if _, ok := agent.Options["path"]; ok {
			file := agent.Options["path"].(string)

			if _, err := os.Stat(filepath.Join(pwd, file)); err == nil {
				file = filepath.Join(pwd, file)
			}
			var err error
			content, err = ioutil.ReadFile(file)
			if err != nil {
				log.Printf(`Error while reading "%s" [%s]`, file, err)
				log.Fatalln(err)
			}
			ucwd = filepath.Dir(file)
		} else if v, ok := agent.Options["url"]; ok {
			response, err := http.Get(v.(string))
			if err != nil {
				log.Fatal(err)
			} else {
				content, err = ioutil.ReadAll(response.Body)
				response.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
			}
			ucwd, _ = os.Getwd()
		} else {
			log.Fatalln("give me a url or a path to read from use plugin")
		}

		fileConfigAgents, err := parseConfig(agent.Name, content, ucwd, "filter")
		if err != nil {
			log.Fatalln("ERROR while using config ", err.Error())
		}

		firstUsedAgent := &fileConfigAgents[len(fileConfigAgents)-1]
		for _, sourceport := range lastOutPorts {
			inPort := config.Port{AgentName: sourceport.AgentName, PortNumber: sourceport.PortNumber}
			firstUsedAgent.XSources = append(firstUsedAgent.XSources, inPort)
		}

		newOutPorts := []config.Port{
			{AgentName: fileConfigAgents[0].Name, PortNumber: 0},
		}
		return fileConfigAgents, newOutPorts, i
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
		{AgentName: agent.Name, PortNumber: 0},
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
				{AgentName: agent.Name, PortNumber: expressionIndex},
			}

			// construire les plugins associés à l'expression
			// en utilisant le outportA
			for pi := 0; pi < len(when.Plugins); pi++ {
				p := when.Plugins[pi]
				var agents []config.Agent
				i++
				// récupérer le dernier outport du plugin créé il devient outportA
				agents, expressionOutPorts, i = buildFilterAgents(logfanname, p, expressionOutPorts, i, pwd)
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
				{AgentName: agent.Name, PortNumber: len(agent.Options["expressions"].(map[int]string)) - 1},
			}
			newOutPorts = append(elseOutPorts, newOutPorts...)
		}
	}

	// ajoute l'agent à la liste des agents construits
	agent_list = append([]config.Agent{agent}, agent_list...)
	return agent_list, newOutPorts, i
}

func isInSlice(needle string, candidates []string) bool {
	for _, symbolType := range candidates {
		if needle == symbolType {
			return true
		}
	}
	return false
}
