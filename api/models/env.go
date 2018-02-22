package models

import "time"

type Env struct {
	Uuid   string `json:"uuid"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Secret bool   `json:"secret"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
