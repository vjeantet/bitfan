package testutils

import cache "github.com/patrickmn/go-cache"

type memory struct {
}

type memorySpace struct {
	name string
	c    *cache.Cache
}

func NewMemory(location string) *memory {
	return &memory{}
}

func (s *memory) close() {
	// persist memory on disk ?
}

func (s *memory) Space(name string) *memorySpace {
	return &memorySpace{
		name: name,
		c:    cache.New(cache.NoExpiration, 0),
	}
}

// Set add an item to the cache, replacing any existing item with the same name
func (m *memorySpace) Set(name string, value interface{}) {
	m.c.Set(name, value, cache.NoExpiration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
func (m *memorySpace) Get(name string) (interface{}, bool) {
	return m.c.Get(name)
}

func (m *memorySpace) Delete(name string) {
	m.c.Delete(name)
}

func (m *memorySpace) IncrementInt(k string, n int) {
	m.c.IncrementInt(k, n)
}

func (m *memorySpace) Items() map[string]interface{} {
	values := m.c.Items()
	r := map[string]interface{}{}
	for key, value := range values {
		r[key] = value.Object
	}
	return r
}
