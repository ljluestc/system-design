package utils

import (
    "net/http"
    "sync"
)

type ProxyTransport struct {
    proxies []string
    index   int
    mu      sync.Mutex
}

func NewProxyTransport(proxies []string) *http.Transport {
    return &http.Transport{} // Placeholder
}

func (t *ProxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
    return http.DefaultTransport.RoundTrip(req)
}