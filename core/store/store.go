package store

import (
	"io"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/timshannon/bolthold"
)

type Store struct {
	db  *bolthold.Store
	log Logger
}

func New(location string, log Logger) (*Store, error) {
	database, err := bolthold.Open(filepath.Join(location, "bitfan.bolt.db"), 0666, nil)
	return &Store{db: database, log: log}, err
}

func (s *Store) Close() {
	s.db.Close()
}

// CopyTo writes the raw database's content to given io.Writer
func (s *Store) CopyTo(w io.Writer) (int, error) {
	size := 0
	err := s.db.Bolt().View(func(tx *bolt.Tx) error {
		size = int(tx.Size())
		_, err := tx.WriteTo(w)
		return err
	})
	return size, err
}
