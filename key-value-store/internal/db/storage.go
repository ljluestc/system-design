package db

import "errors"

type Storage interface {
    Put(key, value string) error
    Get(key string) (string, error)
}

func NewMemoryStorage() Storage {
    return &MemoryStorage{data: make(map[string]string)}
}

type MemoryStorage struct {
    data map[string]string
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