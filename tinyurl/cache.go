package main

import (
    "container/list"
    "sync"
)

// LRUCache is a simple LRU cache
type LRUCache struct {
    capacity int
    items    map[string]*list.Element
    order    *list.List
    mu       sync.RWMutex
}

// Entry is a cache entry
type Entry struct {
    key   string
    value string
}

// NewLRUCache initializes an LRU cache
func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        order:    list.New(),
    }
}

// Set adds or updates a key-value pair
func (c *LRUCache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()

    if elem, ok := c.items[key]; ok {
        c.order.MoveToFront(elem)
        elem.Value.(*Entry).value = value
        return
    }

    if c.order.Len() >= c.capacity {
        oldest := c.order.Back()
        if oldest != nil {
            c.order.Remove(oldest)
            delete(c.items, oldest.Value.(*Entry).key)
        }
    }

    elem := c.order.PushFront(&Entry{key: key, value: value})
    c.items[key] = elem
}

// Get retrieves a value by key
func (c *LRUCache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()

    elem, ok := c.items[key]
    if !ok {
        return "", false
    }
    c.order.MoveToFront(elem)
    return elem.Value.(*Entry).value, true
}