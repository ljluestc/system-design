package cache

import (
    "container/list"
    "sync"
)

type LRUEviction struct {
    mu       sync.Mutex      // Mutex for thread safety
    items    map[string]*list.Element // Key to list element mapping
    list     *list.List      // Doubly-linked list for LRU
    capacity int             // Maximum capacity
}

func NewLRUEviction(capacity int) *LRUEviction {
    return &LRUEviction{
        items:    make(map[string]*list.Element),
        list:     list.New(),
        capacity: capacity,
    }
}

func (l *LRUEviction) Touch(key string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    if elem, ok := l.items[key]; ok {
        l.list.MoveToFront(elem)
    } else {
        if l.list.Len() >= l.capacity {
            oldest := l.list.Back()
            if oldest != nil {
                l.list.Remove(oldest)
                delete(l.items, oldest.Value.(string))
            }
        }
        elem := l.list.PushFront(key)
        l.items[key] = elem
    }
}

func (l *LRUEviction) Remove(key string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    if elem, ok := l.items[key]; ok {
        l.list.Remove(elem)
        delete(l.items, key)
    }
}