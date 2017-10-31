package models

import "time"

type Asset struct {
	ID int `gorm:"primary_key"`

	Uuid string `json:"uuid"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PipelineUUID string

	Name        string
	Type        string
	ContentType string
	Value       []byte //`json:"-"`
	Size        int

	// Base64 encoded content
	// Content string `json:"-" gorm:"-"`
}

const (
	ASSET_TYPE_ENTRYPOINT = "entrypoint"
)
