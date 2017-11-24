package memory

import cache "github.com/patrickmn/go-cache"

type Memory struct {
}

type MemorySpace struct {
	name string
	c    *cache.Cache
}

func New() *Memory {
	return &Memory{}
}

func (s *Memory) Close() {
	// persist Memory ?
}

func (s *Memory) Space(name string) *MemorySpace {
	return &MemorySpace{
		name: name,
		c:    cache.New(cache.NoExpiration, 0),
	}
}

// Set add an item to the cache, replacing any existing item with the same name
func (m *MemorySpace) Set(name string, value interface{}) {
	m.c.Set(name, value, cache.NoExpiration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating whether the key was found.
func (m *MemorySpace) Get(name string) (interface{}, bool) {
	return m.c.Get(name)
}

func (m *MemorySpace) Delete(name string) {
	m.c.Delete(name)
}

func (m *MemorySpace) IncrementInt(k string, n int) {
	m.c.IncrementInt(k, n)
}

func (m *MemorySpace) Items() map[string]interface{} {
	values := m.c.Items()
	r := map[string]interface{}{}
	for key, value := range values {
		r[key] = value.Object
	}
	return r
}
