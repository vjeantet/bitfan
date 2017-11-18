package store

import (
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/timshannon/bolthold"
)

type Store struct {
	db  *bolthold.Store
	log *logrus.Logger
}

func NewStore(location string, log *logrus.Logger) (*Store, error) {
	database, err := bolthold.Open(filepath.Join(location, "bitfan.bolt.db"), 0666, nil)
	return &Store{db: database, log: log}, err
}

func (s *Store) Close() {
	s.db.Close()
}
