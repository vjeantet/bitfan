package models

import "time"

type XProcessor struct {
	Uuid        string   `json:"uuid"`
	Label       string   `json:"label"`
	Behavior    string   `json:"behavior"`
	Stream      bool     `json:"stream"`
	Kind        string   `json:"kind"`
	Args        []string `json:"args"`
	Command     string   `json:"command"`
	Code        string   `json:"code"`
	StdinAs     string   `json:"stdin_as" mapstructure:"stdin_as"`
	StdoutAs    string   `json:"stdout_as" mapstructure:"stdout_as"`
	Description string   `json:"description"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
