package storage

type BlobStore struct {
    data map[string][]byte
}

func NewBlobStore() *BlobStore {
    return &BlobStore{data: make(map[string][]byte)}
}

func (bs *BlobStore) Save(key string, data []byte) {
    bs.data[key] = data
}