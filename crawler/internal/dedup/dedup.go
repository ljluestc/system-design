package dedup

import (
    "crawler/pkg/db"
    "crawler/pkg/models"
)

type Dedup struct {
    db *db.Postgres
}

func New(db *db.Postgres) *Dedup {
    return &Dedup{db: db}
}

func (d *Dedup) IsDuplicateURL(url string) (bool, error) {
    checksum := ComputeChecksum(url)
    return d.db.HasURLChecksum(checksum)
}

func (d *Dedup) SaveURL(url string) error {
    checksum := ComputeChecksum(url)
    return d.db.SaveURLChecksum(checksum)
}

func (d *Dedup) IsDuplicateDocument(doc models.Document) (bool, error) {
    checksum := ComputeChecksum(doc.Content)
    return d.db.HasDocChecksum(checksum)
}

func (d *Dedup) SaveDocument(doc models.Document) error {
    checksum := ComputeChecksum(doc.Content)
    return d.db.SaveDocChecksum(checksum)
}