package cache

import (
	"sync"
	"time"
)

// InMemoryCache - implementation of Cache interface
// local in-memory storage, replacement for Redis in tests
// Mutex is used to protect map (sync.Map can be used too)
type InMemoryCache struct {
	dataMap map[int]inMemoryValue
	lock    *sync.Mutex
}

// inMemoryValue - internal struct to check expiration on values in cache
type inMemoryValue struct {
	SetTime    int64
	Expiration int64
}

// NewInMemoryCache - create new instance of InMemoryCache
func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		dataMap: make(map[int]inMemoryValue, 0),
		lock:    &sync.Mutex{},
	}
}

// Add - add rand value with expiration (in seconds) to cache
func (c *InMemoryCache) Add(key int, expiration int64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dataMap[key] = inMemoryValue{
		SetTime:    time.Now().Unix(),
		Expiration: expiration,
	}

	return nil
}

// Get - check existence of int key in cache
func (c *InMemoryCache) Get(key int) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.dataMap[key]
	if ok && time.Now().Unix()-value.SetTime > value.Expiration {
		return false, nil
	}

	return ok, nil
}

// Delete - delete key from cache
func (c *InMemoryCache) Delete(key int) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.dataMap, key)
}
