package worker

import (
    "crawler/internal/dedup"
    "crawler/internal/dns"
    "crawler/internal/extractor"
    "crawler/internal/fetcher"
    "crawler/internal/scheduler"
    "crawler/internal/storage"
    "crawler/internal/trap"
    "crawler/pkg/models"
)

type Worker struct {
    id         int
    scheduler  *scheduler.Scheduler
    dns        *dns.Resolver
    fetcher    *fetcher.HTTPFetcher
    extractor  *extractor.Extractor
    dedup      *dedup.Dedup
    storage    *storage.BlobStore
    trap       *trap.Detector
}

func New(id int, s *scheduler.Scheduler, d *dns.Resolver, f *fetcher.HTTPFetcher, e *extractor.Extractor, dd *dedup.Dedup, st *storage.BlobStore) *Worker {
    return &Worker{
        id:        id,
        scheduler: s,
        dns:       d,
        fetcher:   f,
        extractor: e,
        dedup:     dd,
        storage:   st,
        trap:      trap.New(),
    }
}

func (w *Worker) Run() {
    for {
        url, err := w.scheduler.NextURL()
        if err != nil {
            continue
        }
        if w.trap.IsTrap(url.Address) {
            continue
        }
        ip, err := w.dns.Resolve(url.Hostname())
        if err != nil {
            continue
        }
        content, err := w.fetcher.Fetch(url.Address)
        if err != nil {
            continue
        }
        doc, newURLs, err := w.extractor.Extract(content)
        if err != nil {
            continue
        }
        if ok, _ := w.dedup.IsDuplicateDocument(doc); ok {
            continue
        }
        if err := w.storage.Save(doc); err != nil {
            continue
        }
        for _, newURL := range newURLs {
            if ok, _ := w.dedup.IsDuplicateURL(newURL); !ok {
                w.scheduler.AddURL(models.URL{Address: newURL, Priority: 1})
                w.dedup.SaveURL(newURL)
            }
        }
    }
}