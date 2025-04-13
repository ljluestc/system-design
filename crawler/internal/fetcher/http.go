package fetcher

import (
    "net/http"
    "time"
    "crawler/internal/utils"
)

type HTTPFetcher struct {
    client *http.Client
}

func NewHTTPFetcher(proxyList []string) *HTTPFetcher {
    return &HTTPFetcher{
        client: &http.Client{
            Timeout: 10 * time.Second,
            Transport: utils.NewProxyTransport(proxyList),
        },
    }
}

func (f *HTTPFetcher) Fetch(url string) ([]byte, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("User-Agent", "CrawlerBot/1.0")

    resp, err := f.client.Do(req)
    if err != nil {
        return nil, utils.Errorf("fetch %s: %v", url, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, utils.Errorf("status %d for %s", resp.StatusCode, url)
    }

    return utils.ReadBody(resp.Body)
}