package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged
	OnEvited func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvited func(string, Value)) *Cache {
	return &Cache{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
		OnEvited: onEvited,
	}
}

// Get function
func (c *Cache) Get(key string) (value Value, ok bool) {
	if item, ok := c.cache[key]; ok {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	item := c.ll.Back()
	if item != nil {
		c.ll.Remove(item)
		kv := item.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes = int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvited != nil {
			c.OnEvited(kv.key, kv.value)
		}
	}
}

// add a value to the cache
func (c *Cache) Add(key string, value Value) {
	if item, ok := c.cache[key]; ok {
		c.ll.MoveToFront(item)
		kv := item.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		item := c.ll.PushFront(&entry{key, value})
		c.cache[key] = item
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
