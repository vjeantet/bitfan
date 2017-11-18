package store

import "time"

type processorSpace struct {
	Key    string `boltholdIndex:"ProcessorKey"`
	Bucket string `boltholdIndex:"ProcessorBucket"`
	Value  []byte

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
