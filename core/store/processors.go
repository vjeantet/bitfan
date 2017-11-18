package store

import "time"

type processor struct {
	Key    string `boltholdIndex:"ProcessorKey"`
	Bucket string `boltholdIndex:"ProcessorBucket"`
	Value  []byte

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *processorStorage) Get(key, bucket string) ([]byte, error) {
	return []byte(""), nil
}

func (s *processorStorage) Set(key, bucket string, value []byte) error {
	return nil
}

func (s *processorStorage) Delete(key, bucket string) error {
	return nil
}

func (s *processorStorage) Has(key, bucket string) (bool, error) {
	return false, nil
}
