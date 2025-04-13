package utils

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    MessagesSent = promauto.NewCounter(prometheus.CounterOpts{
        Name: "messenger_messages_sent_total",
        Help: "Total number of messages sent",
    })
)

func RecordMessageSent() {
    MessagesSent.Inc()
}