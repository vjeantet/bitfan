package runtime

import "github.com/veino/config"

type AgentInformation struct {
	Config config.Agent
	Hooks  map[string]*hook
}
