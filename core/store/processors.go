package store

type processorStorage struct {
	store     *Store
	storeName []byte
}

func (s *Store) NewProcessorStorage(processorType string) (*processorStorage, error) {
	tx, err := s.db.Bolt().Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ps := &processorStorage{
		store:     s,
		storeName: []byte("store_" + processorType),
	}

	_, err = tx.CreateBucketIfNotExists(ps.storeName)
	if err != nil {
		return nil, err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return ps, nil
}

func (p *processorStorage) Get(key, ns string) ([]byte, error) {
	tx, err := p.store.db.Bolt().Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// TRANSACTION
	if ns == "" {
		ns = "default"
	}

	bs := tx.Bucket(p.storeName)
	b := bs.Bucket([]byte(ns))
	if b != nil {
		return b.Get([]byte(key)), nil
	} else {
		return nil, nil
	}

}

func (p *processorStorage) Set(key, ns string, value []byte) error {
	// Start the transaction.
	tx, err := p.store.db.Bolt().Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TRANSACTION
	if ns == "" {
		ns = "default"
	}
	bkt, err := tx.Bucket(p.storeName).CreateBucketIfNotExists([]byte(ns))
	if err != nil {
		return err
	}
	err = bkt.Put([]byte(key), value)
	if err != nil {
		return err
	}
	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *processorStorage) Delete(key, ns string) error {
	// Start the transaction.
	tx, err := p.store.db.Bolt().Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TRANSACTION
	if ns == "" {
		ns = "default"
	}

	bs := tx.Bucket(p.storeName)
	b := bs.Bucket([]byte(ns))
	if b != nil {
		if err := b.Delete([]byte(key)); err != nil {
			return err
		}
		// Commit the transaction.
		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func (p *processorStorage) Has(key, ns string) (bool, error) {
	v, err := p.Get(key, ns)
	if err != nil {
		return false, err
	}
	if v == nil {
		return false, nil
	}
	return true, nil
}
