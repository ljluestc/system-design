package context

import "sync"

type Store struct {
    data  map[string]string
    mutex sync.RWMutex
}

func NewStore() *Store {
    return &Store{data: make(map[string]string)}
}

func (s *Store) SaveHistory(userID, history string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    s.data[userID] = history
}

func (s *Store) GetHistory(userID string) (string, bool) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    h, ok := s.data[userID]
    return h, ok
}