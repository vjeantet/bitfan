package core

import (
	"time"

	"github.com/vjeantet/bitfan/core/config"
)

type Pipeline struct {
	Uuid               string
	Label              string
	agents             map[int]*agent
	ConfigLocation     string
	ConfigHostLocation string
	StartedAt          time.Time
}

func newPipeline(conf *config.Pipeline, configAgents []config.Agent) (*Pipeline, error) {
	p := &Pipeline{
		Uuid:               conf.Uuid,
		Label:              conf.Name,
		ConfigLocation:     conf.ConfigLocation,
		ConfigHostLocation: conf.ConfigHostLocation,
		agents:             map[int]*agent{},
	}

	//normalize
	configAgents = config.Normalize(configAgents)

	// for each agents in configAgents (outputs first)
	orderedAgentConfList := config.Sort(configAgents, config.SortInputsFirst)
	for _, agentConf := range orderedAgentConfList {
		agentConf.PipelineUUID = p.Uuid
		agentConf.PipelineName = p.Label
		a, err := newAgent(agentConf)
		if err != nil {
			Log().Errorf("%s agent '%-d' can not start", agentConf.Type, agentConf.ID)
			return nil, err
		}

		// register input chan for futur reference and connecting
		// for each sources
		for _, sourcePort := range agentConf.AgentSources {
			// find agent source.ID aSource
			aSource := p.agents[sourcePort.AgentID]
			// add a(in) to aSource outputs with port
			aSource.addOutput(a.packetChan, sourcePort.PortNumber)
		}
		p.addAgent(a)
	}

	return p, nil
}

func (p *Pipeline) addAgent(a *agent) error {
	a.conf.PipelineName = p.Label
	a.conf.PipelineUUID = p.Uuid
	p.agents[a.ID] = a

	return nil
}

// Start all agents, begin with last
func (p *Pipeline) start() error {
	orderedAgentConfList := config.Sort(p.agentsConfiguration(), config.SortOutputsFirst)
	for _, agentConf := range orderedAgentConfList {
		Log().Debugf("start %d - %s", agentConf.ID, p.agents[agentConf.ID].Label)
		p.agents[agentConf.ID].start()
	}
	p.StartedAt = time.Now()
	return nil
}

func (p *Pipeline) agentsConfiguration() []config.Agent {
	agentsConf := []config.Agent{}
	for _, a := range p.agents {
		agentsConf = append(agentsConf, a.conf)
	}
	return agentsConf
}

// Stop all agents, begin with first
func (p *Pipeline) stop() error {
	orderedAgentConfList := config.Sort(p.agentsConfiguration(), config.SortInputsFirst)
	for _, agentConf := range orderedAgentConfList {
		Log().Debugf("stop %d - %s", agentConf.ID, p.agents[agentConf.ID].Label)
		p.agents[agentConf.ID].stop()
	}
	return nil
}
