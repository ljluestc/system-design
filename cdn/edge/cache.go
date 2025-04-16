package edge

import (
    "container/list"
)

type Cache struct {
    capacity int
    items    map[string]*list.Element
    list     *list.List
}

type entry struct {
    key, value string
}

func NewCache(capacity int) *Cache {
    return &Cache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (c *Cache) Get(key string) (string, bool) {
    if elem, ok := c.items[key]; ok {
        c.list.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    return "", false
}

func (c *Cache) Put(key, value string) {
    if elem, ok := c.items[key]; ok {
        c.list.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }
    if c.list.Len() >= c.capacity {
        oldest := c.list.Back()
        if oldest != nil {
            delete(c.items, oldest.Value.(*entry).key)
            c.list.Remove(oldest)
        }
    }
    elem := c.list.PushFront(&entry{key, value})
    c.items[key] = elem
}