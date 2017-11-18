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
