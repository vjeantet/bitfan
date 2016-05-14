package config

type Agent struct {
	Pipeline       string
	Name           string   `json:"name"`
	Sources        []string `json:"sources"`
	XSources       PortList
	XRecipients    PortList
	Type           string                 `json:"type"`
	Schedule       string                 `json:"schedule"`
	EventRetention string                 `json:"event_retention"`
	PoolSize       int                    `json:"pool_size"`
	Buffer         int                    `json:"buffer_size"`
	Options        map[string]interface{} `json:"options"`
}

type AgentList map[string]Agent

var agentConfigurations = AgentList{}

type Pipelines map[string]AgentList

// All returns loaded Agents configurations
func Agents() AgentList {
	return agentConfigurations
}

func (a AgentList) NamesSort(sortOrder int) ([]string, []string) {
	var agentsDependencyGraph = graph{}

	for _, agentConfiguration := range a {
		agentsDependencyGraph[agentConfiguration.Name] = func() []string {
			sources := []string{}
			for _, agentname := range agentConfiguration.XSources {
				sources = append(sources, agentname.AgentName)
			}
			return sources
		}()
	}

	order, cycle := topSortDFS(agentsDependencyGraph)

	if sortOrder == SortInputsFirst {
		order = reverseList(order)
	}
	return order, cycle
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
		agentsDependencyGraph[agentConfiguration.Name] = func() []string {
			sources := []string{}
			for _, agentname := range agentConfiguration.XSources {
				sources = append(sources, agentname.AgentName)
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
			if agentConf.Name == n {
				sac = append(sac, agentConf)
			}
		}
	}

	return sac
}

// Add agent's configuration to veino
func AddAgent(agentConf *Agent) {
	addAgent(agentConf)
}

func addAgent(agentConf *Agent) {
	// Find recipients
	Normalize(agentConf)

	// Find Existing Agent Conf for each Source, and add agentConf.Name to its Recipients
	for _, agentPort := range agentConf.XSources {
		if ac, ok := agentConfigurations[agentPort.AgentName]; ok {
			ac.XRecipients = append(ac.XRecipients, Port{AgentName: agentConf.Name, PortNumber: agentPort.PortNumber})
			agentConfigurations[agentPort.AgentName] = ac
		}
	}
	// fmt.Printf("    %-10s %s\n", agentConf.Name, agentConf.XRecipients.String())
	agentConfigurations[agentConf.Name] = *agentConf
}

// RemoveAgent delete agent configuration from veino
func RemoveAgent(agentName string) {
	// Find all source agent of agentName
	for _, srcAgentPort := range agentConfigurations[agentName].XSources {
		tmp := agentConfigurations[srcAgentPort.AgentName]
		for k, v := range tmp.XRecipients {
			if v.AgentName == agentName {
				tmp.XRecipients = append(tmp.XRecipients[:k], tmp.XRecipients[k+1:]...)
				agentConfigurations[srcAgentPort.AgentName] = tmp
				break
			}
		}
	}

	delete(agentConfigurations, agentName)
}

// FindRecipientNames returns recipients agents names having source : sourceAgentName[sourceAgentPort]
func FindRecipientNames(agentName string, agentPort int) []string {
	r := []string{}
	for _, aport := range agentConfigurations[agentName].XRecipients {
		if aport.PortNumber == agentPort {
			r = append(r, aport.AgentName)
		}
	}
	return r
}

// WhoWaitForThisAgentName returns agents recipients as portList
func Normalize(agentConf *Agent) {
	agentConf.XRecipients = whoWaitForThisAgentName(agentConf.Name)
}

func whoWaitForThisAgentName(agentName string) PortList {
	var recipentAgents = PortList{}

	for _, agentConfiguration := range agentConfigurations {
		for _, sourceAgentPort := range agentConfiguration.XSources {
			if sourceAgentPort.AgentName == agentName { // ok on a un agent qui se source sur agentName
				for _, v := range agentConfiguration.XSources {
					if agentName != v.AgentName {
						continue
					}
					recipentAgents = append(recipentAgents, Port{AgentName: agentConfiguration.Name, PortNumber: v.PortNumber})
				}
			}
		}
	}

	return recipentAgents
}

const (
	SortInputsFirst = iota + 1
	SortOutputsFirst
)

func reverseList(s []string) (r []string) {
	for _, i := range s {
		i := i
		defer func() { r = append(r, i) }()
	}
	return
}

type graph map[string][]string

func topSortDFS(g graph) (order, cyclic []string) {
	L := make([]string, len(g))
	i := len(L)
	temp := map[string]bool{}
	perm := map[string]bool{}
	var cycleFound bool
	var cycleStart string
	var visit func(string)
	visit = func(n string) {
		switch {
		case temp[n]:
			cycleFound = true
			cycleStart = n
			return
		case perm[n]:
			return
		}
		temp[n] = true
		for _, m := range g[n] {
			visit(m)
			if cycleFound {
				if cycleStart > "" {
					cyclic = append(cyclic, n)
					if n == cycleStart {
						cycleStart = ""
					}
				}
				return
			}
		}
		delete(temp, n)
		perm[n] = true
		i--
		L[i] = n
	}
	for n := range g {
		if perm[n] {
			continue
		}
		visit(n)
		if cycleFound {
			return nil, cyclic
		}
	}
	return L, nil
}
