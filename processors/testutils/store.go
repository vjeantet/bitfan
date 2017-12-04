package testutils

type processorStorage struct {
	storeName []byte
}

func newStore(p *DummyProcessorContext) *processorStorage {
	return &processorStorage{}
}

func (p *processorStorage) Get(key, ns string) ([]byte, error) {
	return nil, nil
}

func (p *processorStorage) Set(key, ns string, value []byte) error {
	return nil
}

func (p *processorStorage) Delete(key, ns string) error {
	return nil
}

func (p *processorStorage) Has(key, ns string) (bool, error) {
	return true, nil
}
