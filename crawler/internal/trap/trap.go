package trap

import (
    "strings"
    "net/url"
)

type Detector struct {
    seen map[string]int
}

func New() *Detector {
    return &Detector{seen: make(map[string]int)}
}

func (d *Detector) IsTrap(urlStr string) bool {
    u, err := url.Parse(urlStr)
    if err != nil {
        return true
    }
    path := u.Path
    if strings.Contains(path, "../") || strings.Count(path, "/") > 10 {
        return true
    }
    d.seen[urlStr]++
    return d.seen[urlStr] > 5
}