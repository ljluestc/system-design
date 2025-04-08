package integration

import (
    "net/http"
    "testing"
)

func TestShardingInsert(t *testing.T) {
    resp, err := http.Get("http://localhost:3001/api/sharding/get?userId=user1")
    if err != nil || resp.StatusCode != 200 {
        t.Errorf("API failed: %v", err)
    }
}