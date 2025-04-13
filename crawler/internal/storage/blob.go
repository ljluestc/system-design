package storage

import (
    "crawler/internal/utils"
    "crawler/pkg/models"
)

type BlobStore struct {
    bucket string
}

func NewBlobStore(bucket string) *BlobStore {
    return &BlobStore{bucket: bucket}
}

func (s *BlobStore) Save(doc models.Document) error {
    return utils.UploadToS3(s.bucket, doc.ID, []byte(doc.Content))
}