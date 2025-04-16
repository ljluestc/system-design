package edge

import (
    "cdn/origin"
    "sync"
)

type Server struct {
    ID        int
    Cache     *Cache
    Origin    *origin.Server
    Compressor *Compressor
    mutex     sync.RWMutex
}

func NewServer(id, capacity int) *Server {
    return &Server{
        ID:         id,
        Cache:      NewCache(capacity),
        Origin:     origin.NewServer(),
        Compressor: NewCompressor(),
    }
}

func (s *Server) Get(key string) (string, bool) {
    s.mutex.RLock()
    if value, ok := s.Cache.Get(key); ok {
        s.mutex.RUnlock()
        return s.Compressor.Compress(value), true
    }
    s.mutex.RUnlock()
    value, ok := s.Origin.Fetch(key)
    if ok {
        s.mutex.Lock()
        s.Cache.Put(key, value)
        s.mutex.Unlock()
        return s.Compressor.Compress(value), true
    }
    return "", false
}

func (s *Server) Put(key, value string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.Cache.Put(key, value)
}