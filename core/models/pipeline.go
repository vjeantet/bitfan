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
	ID int `json:"-" gorm:"primary_key"`

	Uuid string `json:"uuid"`

	// the Label
	Label string `json:"label"`

	Description string

	// the location
	ConfigLocation string `json:"config_location" gorm:"-"`

	// the location's host
	ConfigHostLocation string `json:"config_host_location" gorm:"-"`

	// Assets
	Assets []Asset `json:"assets" gorm:"ForeignKey:PipelineUUID;AssociationForeignKey:Uuid"`

	Active       bool `gorm:"-"`
	LocationPath string

	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt time.Time `json:"started_at"`
}
