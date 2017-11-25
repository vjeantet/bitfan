package core

// Sort will return a sorted list of config.Agent,
// it sorts agents by computing links dependencies between them
//
// use sortOrder param config.SortInputsFirst to get agents which are not waiting events (no sources) firstly (like inputs)
//
// use sortOrder param config.SortOutputsFirst to get agents which are not sources of any other agents firstly (like outputs)
func Sort(agentConflist map[int]*Agent, sortOrder int) []*Agent {
	sac := []*Agent{}
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

// WhoWaitForThisAgentName returns agents recipients as portList
func whoWaitForThisAgentID(ID int, agentConfigurations map[int]*Agent) PortList {
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
