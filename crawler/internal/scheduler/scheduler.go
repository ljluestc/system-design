package scheduler

import (
    "crawler/pkg/db"
    "crawler/pkg/models"
    "crawler/internal/storage"
)

type Scheduler struct {
    db    *db.Postgres
    redis *storage.Redis
    queue *PriorityQueue
}

func New(db *db.Postgres, redis *storage.Redis) *Scheduler {
    return &Scheduler{
        db:    db,
        redis: redis,
        queue: NewPriorityQueue(),
    }
}

func (s *Scheduler) AddURL(url models.URL) error {
    // Save URL to Postgres for persistence
    if err := s.db.SaveURL(url); err != nil {
        return err
    }
    // Push to priority queue
    return s.queue.Push(url)
}

func (s *Scheduler) NextURL() (models.URL, error) {
    // Pop next URL from priority queue
    url, err := s.queue.Pop()
    if err != nil {
        return models.URL{}, err
    }
    return url, nil
}