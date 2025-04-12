package core

import (
    "container/list"
)

// LRUCache is a basic Least Recently Used cache
type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
}

// Entry holds a cache key-value pair
type Entry struct {
    key   string
    value []string
}

// NewLRUCache creates a new LRUCache
func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

// Get retrieves a value from the cache
func (c *LRUCache) Get(key string) ([]string, bool) {
    if elem, exists := c.cache[key]; exists {
        c.list.MoveToFront(elem)
        return elem.Value.(*Entry).value, true
    }
    return nil, false
}

// Put adds a value to the cache
func (c *LRUCache) Put(key string, value []string) {
    if elem, exists := c.cache[key]; exists {
        c.list.MoveToFront(elem)
        elem.Value.(*Entry).value = value
        return
    }
    if c.list.Len() >= c.capacity {
        oldest := c.list.Back()
        if oldest != nil {
            c.list.Remove(oldest)
            delete(c.cache, oldest.Value.(*Entry).key)
        }
    }
    entry := &Entry{key: key, value: value}
    elem := c.list.PushFront(entry)
    c.cache[key] = elem
}