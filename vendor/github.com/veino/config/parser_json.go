package config

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

func ParseJSON(content []byte) ([]Agent, error) {
	agentConfList := []Agent{}
	dec := json.NewDecoder(bytes.NewReader(content))

	for {
		var conf Agent
		if err := dec.Decode(&conf); err == io.EOF {
			break
		} else if err != nil {
			//TODO: handle error
			return nil, err
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

		agentConfList = append(agentConfList, conf)
	}

	return agentConfList, nil
}
