package scheduler

import (
    "container/heap"
    "crawler/pkg/models"
    "errors"
)

type PriorityQueue struct {
    urls []models.URL
}

func NewPriorityQueue() *PriorityQueue {
    pq := &PriorityQueue{}
    heap.Init(pq)
    return pq
}

func (pq *PriorityQueue) Push(url models.URL) error {
    heap.Push(pq, url)
    return nil
}

func (pq *PriorityQueue) Pop() (models.URL, error) {
    if pq.Len() == 0 {
        return models.URL{}, errors.New("queue empty")
    }
    return heap.Pop(pq).(models.URL), nil
}

func (pq *PriorityQueue) Len() int {
    return len(pq.urls)
}

func (pq *PriorityQueue) Less(i, j int) bool {
    return pq.urls[i].Priority > pq.urls[j].Priority // Higher priority first
}

func (pq *PriorityQueue) Swap(i, j int) {
    pq.urls[i], pq.urls[j] = pq.urls[j], pq.urls[i]
}

func (pq *PriorityQueue) PushItem(x interface{}) {
    pq.urls = append(pq.urls, x.(models.URL))
}

func (pq *PriorityQueue) PopItem() interface{} {
    n := len(pq.urls)
    item := pq.urls[n-1]
    pq.urls = pq.urls[:n-1]
    return item
}