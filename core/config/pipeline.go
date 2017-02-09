package config

import (
	"time"

	fqdn "github.com/ShowMax/go-fqdn"
)

type PipelineState int

type Pipeline struct {
	ID                 int
	Name               string
	Description        string
	ConfigLocation     string
	ConfigHostLocation string

	StartedAt time.Time
	StoppedAt time.Time
}

var pipelineIndex int = 0

func NewPipeline(name, description, configLocation string) *Pipeline {
	pipelineIndex++
	return &Pipeline{
		ID:                 pipelineIndex,
		Name:               name,
		Description:        description,
		ConfigLocation:     configLocation,
		ConfigHostLocation: fqdn.Get(),
	}
}
