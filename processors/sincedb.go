package processors

import (
	"sync"
	"time"
)

type SinceDB struct {
	options *SinceDBOptions

	sinceDBLastInfosRaw []byte
	sinceDBLastSaveTime time.Time
	sinceDBInfosMutex   *sync.Mutex
}

type SinceDBOptions struct {
	Identifier    string
	WriteInterval int
}

// NewSinceDB loadExisting data from datastore according to the Identifier option.
func NewSinceDB(sdboptions *SinceDBOptions) *SinceDB {
	s := &SinceDB{
		options: sdboptions,
	}

	return s
}

// Save SinceDB Offsets to Storage
func (s *SinceDB) Save() error {
	return nil
}

// Retreive SinceDB ressource's offset from Storage
func (s *SinceDB) RessourceOffset(ressource string) (int, error) {
	return 0, nil
}

// Update a ressource's offset
func (s *SinceDB) SetRessourceOffset(ressource string, offset int) error {

	// Remember ressource offset to write

	// Check LastTimeWeSave

	// If to soon
	// -- Assur a Ticker
	// -- Return

	// If it's time
	// -- Save()

	return nil
}
