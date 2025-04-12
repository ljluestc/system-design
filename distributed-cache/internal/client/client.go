package client

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

type CacheClient struct {
    baseURL string
}

func NewCacheClient(baseURL string) *CacheClient {
    return &CacheClient{baseURL: baseURL}
}

func (c *CacheClient) Get(key string) (string, error) {
    resp, err := http.Get(c.baseURL + "/cache/" + key)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}

func (c *CacheClient) Set(key, value string) error {
    req, err := http.NewRequest("PUT", c.baseURL+"/cache/"+key, strings.NewReader(value))
    if err != nil {
        return err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    return nil
}