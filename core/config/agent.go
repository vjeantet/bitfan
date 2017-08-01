package config

type Agent struct {
	ID              int      `json:"id"`
	Label           string   `json:"label"`
	Sources         []string `json:"sources"`
	PipelineName    string
	PipelineID      int
	AgentSources    PortList
	AgentRecipients PortList
	Type            string                 `json:"type"`
	Schedule        string                 `json:"schedule"`
	PoolSize        int                    `json:"pool_size"`
	Buffer          int                    `json:"buffer_size"`
	Options         map[string]interface{} `json:"options"`
	Wd              string
}

var agentIndex int = 0

func NewAgent() Agent {
	agentIndex++
	return Agent{
		ID: agentIndex,
	}
}

// Sort will return a sorted list of config.Agent,
// it sorts agents by computing links dependencies between them
//
// use sortOrder param config.SortInputsFirst to get agents which are not waiting events (no sources) firstly (like inputs)
//
// use sortOrder param config.SortOutputsFirst to get agents which are not sources of any other agents firstly (like outputs)
func Sort(agentConflist []Agent, sortOrder int) []Agent {
	sac := []Agent{}
	// sac = append(sac, agentConflist...)

	var agentsDependencyGraph = graph{}

	for _, agentConfiguration := range agentConflist {
		agentsDependencyGraph[agentConfiguration.ID] = func() []int {
			sources := []int{}
			for _, port := range agentConfiguration.AgentSources {
				sources = append(sources, port.AgentID)
			}
			return sources
		}()
	}

	order, _ := topSortDFS(agentsDependencyGraph)

	if sortOrder == SortInputsFirst {
		order = reverseList(order)
	}

	for _, n := range order {
		for _, agentConf := range agentConflist {
			if agentConf.ID == n {
				sac = append(sac, agentConf)
			}
		}
	}

	return sac
}

func Normalize(agentConf []Agent) []Agent {
	for k := range agentConf {
		agentConf[k].AgentRecipients = whoWaitForThisAgentID(agentConf[k].ID, agentConf)
	}
	return agentConf
}

// WhoWaitForThisAgentName returns agents recipients as portList
func whoWaitForThisAgentID(ID int, agentConfigurations []Agent) PortList {
	var recipentAgents = PortList{}

	for _, agentConfiguration := range agentConfigurations {
		for _, sourceAgentPort := range agentConfiguration.AgentSources {
			if sourceAgentPort.AgentID == ID { // ok on a un agent qui se source sur ID
				for _, v := range agentConfiguration.AgentSources {
					if ID != v.AgentID {
						continue
					}
					recipentAgents = append(recipentAgents, Port{AgentID: agentConfiguration.ID, PortNumber: v.PortNumber})
				}
			}
		}
	}

	return recipentAgents
}
