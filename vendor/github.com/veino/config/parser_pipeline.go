package config

import (
	"log"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/clbanning/mxj"
)

func ParsePIPELINE(pipelineName string, content []byte) ([]Agent, error) {
	agentConfList := []Agent{}

	var pipelineMap map[string]interface{}
	md, err := toml.Decode(string(content), &pipelineMap)
	if err != nil {
		log.Fatal(err)
	}

	lastAgentName := ""
	i := 0
	for _, key := range md.Keys() {
		if false == strings.Contains(key.String(), ".") {
			i++
			procKind := key.String()
			procOptions := pipelineMap[key.String()]
			agentName := pipelineName + "-" + strconv.Itoa(i) + "-" + procKind

			var conf Agent
			conf.Type = procKind
			conf.Name = agentName

			if d, err := mxj.Map(procOptions.(map[string]interface{})).ValueForPath("pool"); err == nil {
				conf.PoolSize = int(d.(int64))
			}
			mxj.Map(procOptions.(map[string]interface{})).Remove("pool")

			interval := mxj.Map(procOptions.(map[string]interface{})).ValueOrEmptyForPathString("interval")
			conf.Schedule = interval
			mxj.Map(procOptions.(map[string]interface{})).Remove("interval")
			conf.Options = procOptions.(map[string]interface{})
			if lastAgentName != "" {
				conf.Sources = []string{lastAgentName}
			}

			conf.XSources = PortList{}
			for _, sourceport := range conf.Sources {
				vals := strings.Split(sourceport, "@")
				name := vals[0]
				portNum := 0
				if len(vals) > 1 {
					portNum, _ = strconv.Atoi(vals[1])
				}
				conf.XSources = append(conf.XSources, Port{AgentName: name, PortNumber: portNum})
			}

			if conf.PoolSize == 0 {
				conf.PoolSize = 1
			}

			agentConfList = append([]Agent{conf}, agentConfList...)

			lastAgentName = agentName
		}

	}

	return agentConfList, nil
}
