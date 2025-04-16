package storage

type MemoryStore struct {
    data map[string]string
}

func NewMemoryStore() *MemoryStore {
    return &MemoryStore{data: make(map[string]string)}
}

func (s *MemoryStore) Get(key string) (string, bool) {
    v, ok := s.data[key]
    return v, ok
}

func (s *MemoryStore) Put(key, value string) {
    s.data[key] = value
}