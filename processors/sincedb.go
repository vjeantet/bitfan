package processors

import (
	"strconv"
	"time"

	"golang.org/x/sync/syncmap"
)

type SinceDB struct {
	options *SinceDBOptions
	offsets *syncmap.Map
	done    chan (bool)
	dryrun  bool
}

type SinceDBOptions struct {
	Identifier    string
	WriteInterval int
	Storage       IStore
}

// NewSinceDB loadExisting data from datastore according to the Identifier option.
func NewSinceDB(sdboptions *SinceDBOptions) *SinceDB {
	s := &SinceDB{
		options: sdboptions,
		done:    make(chan (bool)),
		offsets: &syncmap.Map{},
	}

	if s.options.Identifier == "/dev/null" || s.options.Identifier == "" {
		s.dryrun = true
	}

	if s.options.WriteInterval == 0 {
		s.options.WriteInterval = 15
	}

	if s.dryrun == false {
		// Start the write looper
		go func() {
			tick := time.NewTicker(time.Duration(s.options.WriteInterval) * time.Second)
			defer tick.Stop()
			for {
				select {
				case <-tick.C:
					s.save()
				case <-s.done:
					return
				}
			}
		}()
	}

	return s
}

func (s *SinceDB) save() {
	if s.dryrun {
		return
	}
	s.offsets.Range(func(key, value interface{}) bool {
		s.options.Storage.Set(key.(string), s.options.Identifier, value.([]byte))
		s.offsets.Delete(key)
		return true
	})
}

// Save SinceDB Offsets to Storage
func (s *SinceDB) Close() error {
	close(s.done)
	s.save()
	return nil
}

// Update a ressource's offset
func (s *SinceDB) SetRessource(ressource string, v []byte) error {
	s.offsets.Store(ressource, v)
	return nil
}

func (s *SinceDB) Ressource(ressource string) ([]byte, error) {

	// If a value not already stored exists
	if value, ok := s.offsets.Load(ressource); ok {
		return value.([]byte), nil
	}

	// Try to find value in storage
	v, err := s.options.Storage.Get(ressource, s.options.Identifier)
	if err != nil {
		return []byte{}, err
	}

	return v, nil
}

// Retreive SinceDB ressource's offset from Storage
func (s *SinceDB) RessourceOffset(ressource string) (int, error) {

	// If a value not already stored exists
	if value, ok := s.offsets.Load(ressource); ok {
		offset, err := strconv.Atoi(string(value.([]byte)))
		if err != nil {
			return 0, err
		}
		return offset, nil
	}

	// Try to find value in storage
	v, err := s.options.Storage.Get(ressource, s.options.Identifier)
	if err != nil {
		return 0, err
	}

	offset, _ := strconv.Atoi(string(v))
	if err != nil {
		return 0, err
	}

	return offset, nil
}

// Update a ressource's offset
func (s *SinceDB) SetRessourceOffset(ressource string, offset int) error {
	sOffset := strconv.Itoa(offset)
	s.offsets.Store(ressource, []byte(sOffset))
	return nil
}
