package storage

type MemoryStorage struct {
    data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
    return &MemoryStorage{data: make(map[string]string)}
}

func (m *MemoryStorage) Get(key string) (string, bool) {
    val, ok := m.data[key]
    return val, ok
}

func (m *MemoryStorage) Set(key, value string) {
    m.data[key] = value
}