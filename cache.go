package cache

import (
	"container/list"
	"sync"
	"time"
)

// Items represents an individual cache item
type item[T any] struct {
	key      string
	value    T
	exiresAt time.Time
	element  *list.Element
}

// Cache is a thread-safe generic LRU cach with TTL support
type Cache[T any] struct {
	mutex   sync.Mutex
	items   map[string]*item[T]
	order   *list.List
	maxSize int
}

// New creates a new cache with a given max size
func New[T any](maxSize int) *Cache[T] {
	return &Cache[T]{
		items:   make(map[string]*item[T]),
		order:   list.New(),
		maxSize: maxSize,
	}
}

func (c *Cache[T]) Set(key string, value T, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	now := time.Now()
	expiry := now.Add(ttl)

	// update
	if item, found := c.items[key]; found {
		item.value = value
		item.exiresAt = expiry
		c.order.MoveToFront(item.element)
		return
	}

	// evict
	if len(c.items) >= c.maxSize {
		c.evictLRU()
	}

	// add
	elem := c.order.PushFront(key)
	c.items[key] = &item[T]{
		key:      key,
		value:    value,
		exiresAt: expiry,
		element:  elem,
	}
}

func (c *Cache[T]) Get(key string) (T, bool) {
	var z T

	c.mutex.Lock()
	defer c.mutex.Unlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.exiresAt) {
		// found but expired
		if found {
			c.remove(key)
		}
		return z, false
	}

	// move to front
	c.order.MoveToFront(item.element)
	return item.value, true
}

func (c *Cache[T]) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.remove(key)
}

func (c *Cache[T]) evictLRU() {
	if back := c.order.Back(); back != nil {
		key := back.Value.(string)
		c.remove(key)
	}
}

func (c *Cache[T]) remove(key string) {
	if e, ok := c.items[key]; ok {
		c.order.Remove(e.element)
		delete(c.items, key)
	}
}
