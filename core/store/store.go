package store

import (
	"path/filepath"

	"github.com/karlseguin/gerb/core"
	"github.com/timshannon/bolthold"
)

type Store struct {
	db  *bolthold.Store
	log core.Logger
}

func NewStore(location string, log core.Logger) (*Store, error) {
	database, err := bolthold.Open(filepath.Join(location, "bitfan.bolt.db"), 0666, nil)
	return &Store{db: database, log: log}, err
}

func (s *Store) Close() {
	s.db.Close()
}

type processorStorage struct {
	store         *Store
	processorType string
}

func (s *Store) NewProcessorStorage(processorType string) *processorStorage {
	return &processorStorage{store: s, processorType: processorType}
}

func (p *processorStorage) Get(key, bucket string) ([]byte, error) {
	return []byte(""), nil
}

func (p *processorStorage) Set(key, bucket string, value []byte) error {
	return nil
}

func (p *processorStorage) Delete(key, bucket string) error {
	return nil
}

func (p *processorStorage) Has(key, bucket string) (bool, error) {
	return false, nil
}
