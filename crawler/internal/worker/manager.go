package worker

import (
    "crawler/internal/dedup"
    "crawler/internal/dns"
    "crawler/internal/extractor"
    "crawler/internal/fetcher"
    "crawler/internal/scheduler"
    "crawler/internal/storage"
)

type Manager struct {
    workers []*Worker
}

func NewManager(s *scheduler.Scheduler, d *dns.Resolver, f *fetcher.HTTPFetcher, e *extractor.Extractor, dd *dedup.Dedup, st *storage.BlobStore, count int) *Manager {
    m := &Manager{}
    for i := 0; i < count; i++ {
        m.workers = append(m.workers, New(i, s, d, f, e, dd, st))
    }
    return m
}

func (m *Manager) Start() {
    for _, w := range m.workers {
        go w.Run()
    }
}