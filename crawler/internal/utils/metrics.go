package utils

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    URLsCrawled = promauto.NewCounter(prometheus.CounterOpts{
        Name: "crawler_urls_crawled_total",
        Help: "Total URLs crawled",
    })
)

func RecordURLCrawled() {
    URLsCrawled.Inc()
}