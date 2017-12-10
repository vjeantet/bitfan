package models

import "time"

// Pipeline represents a pipeline
//
// A Pipeline is ....
//
// A Pipeline can have.....
//
// swagger:model Pipeline
type Pipeline struct {
	Uuid string `json:"uuid"`

	Playground bool

	// the Label
	Label string `json:"label"`

	Description string

	// the location
	ConfigLocation string `json:"config_location"`

	// the location's host
	ConfigHostLocation string `json:"config_host_location"`

	// Assets
	Assets []Asset `json:"assets"`

	Active       bool
	LocationPath string

	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt time.Time `json:"started_at"`

	AutoStart bool `json:"auto_start" mapstructure:"auto_start"`
}
