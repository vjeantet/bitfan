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

	Content string

	Version int

	Assets []Asset `gorm:"many2many:pipeline_assets;"`
}

type Asset struct {
	ID int `gorm:"primary_key"`

	Name        string
	Type        string
	ContentType string
	Value       []byte
	Size        int
}
