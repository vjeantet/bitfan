package ui

import "time"

// Pipeline represents a pipeline
type Pipeline struct {
	ID int `gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time
	// Static ID
	Uuid string

	// the Label
	Label string

	Description string

	Assets []Asset `gorm:"ForeignKey:PipelineID"`
}

type Asset struct {
	ID int `gorm:"primary_key"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PipelineID int

	Name        string
	Type        string
	ContentType string
	Value       []byte
	Size        int
}
