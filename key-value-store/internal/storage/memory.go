package storage

import "errors"

type MemoryStorage struct {
    data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{data: make(map[string]string)}
}

func (s *MemoryStorage) Put(key, value string) error {
    s.data[key] = value
    return nil
}

func (s *MemoryStorage) Get(key string) (string, error) {
    val, ok := s.data[key]
    if !ok {
        return "", errors.New("key not found")
    }
    return val, nil
}