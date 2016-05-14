package memory

import "github.com/pmylund/go-cache"

func NewMemory(namespace string) *Memory {
	m := &Memory{}
	m.Namespace(namespace)
	return m
}

type Memory struct {
	namespace string
	c         *cache.Cache
}

func (m *Memory) Namespace(ns string) {
	m.namespace = ns
	m.c = cache.New(cache.NoExpiration, 0)
}

// Set add an item to the cache, replacing any existing item with the same name
func (m *Memory) Set(name string, value interface{}) {
	m.c.Set(name, value, cache.NoExpiration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
func (m *Memory) Get(name string) (interface{}, bool) {
	return m.c.Get(name)
}

func (m *Memory) Delete(name string) {
	m.c.Delete(name)
}

func (m *Memory) IncrementInt(k string, n int) {
	m.c.IncrementInt(k, n)
}

func (m *Memory) Items() map[string]interface{} {
	values := m.c.Items()
	r := map[string]interface{}{}
	for key, value := range values {
		r[key] = value.Object
	}
	return r
}
