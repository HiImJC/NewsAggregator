package cache

import (
	"sync"
)

type Cache struct {
	sync.RWMutex
	data map[string]interface{}
}

// New creates a new instance of a working Cache.
func New() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

// Put places the given value within the Cache, for
// retrieval later.
func (c *Cache) Put(k string, v interface{}) {
	c.Lock()
	defer c.Unlock()

	c.data[k] = v
}

// Get returns the cached value for the provided key
// (nil if not present) and a bool representing the presense
// of the key in the Cache.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	cachedValue, ok := c.data[k]
	return cachedValue, ok
}

// Delete removes the key from the Cache. NOP if key is not present
func (c *Cache) Delete(k string) {
	c.Lock()
	defer c.Unlock()

	delete(c.data, k)
}

// Consume retrieves the value from the Cache, then removes the
// item.
// Returns: Cache value (nil if not present), bool representing
// presense of key in the Cache
func (c *Cache) Consume(k string) (interface{}, bool) {
	cacheValue, ok := c.Get(k)
	if ok {
		c.Delete(k)
	}

	return cacheValue, ok
}

// Keys returns back a slice containing all keys present in the cache
func (c *Cache) Keys() []string {
	c.Lock()
	defer c.Unlock()

	keys := make([]string, 0, len(c.data))
	for k, _ := range c.data {
		keys = append(keys, k)
	}

	return keys
}
