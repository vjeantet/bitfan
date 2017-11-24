package models

import "time"

type Asset struct {
	Uuid string `json:"uuid"`

	CreatedAt time.Time
	UpdatedAt time.Time

	PipelineUUID string

	Name        string
	Type        string
	ContentType string
	Value       []byte
	Size        int
}

const (
	ASSET_TYPE_ENTRYPOINT = "entrypoint"
)
