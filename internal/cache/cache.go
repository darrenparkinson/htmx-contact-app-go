package cache

import (
	"sync"
	"time"
)

type Cache[T any, U comparable] struct {
	lastUpdated time.Time
	entries     []T
	entriesMap  map[U]T
	mu          sync.RWMutex
}

// NewCache constructs a new cache containing a list of entries of the provided type
func NewCache[T any, U comparable](entries []T, entriesMap map[U]T) *Cache[T, U] {
	return &Cache[T, U]{
		entries:     entries,
		entriesMap:  entriesMap,
		lastUpdated: time.Now(),
	}
}

func (c *Cache[T, U]) LastUpdated() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUpdated
}

func (c *Cache[T, U]) Update(entries []T, entriesMap map[U]T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = entries
	c.entriesMap = entriesMap
	c.lastUpdated = time.Now()
}

func (c *Cache[T, U]) Retrieve() []T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.entries
}

func (c *Cache[T, U]) RetrieveMap() map[U]T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.entriesMap
}

func (c *Cache[T, U]) RetrieveMapEntry(id U) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entriesMap[id]
	return entry, ok
}
