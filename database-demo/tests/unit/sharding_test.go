package unit

import (
    "testing"
    "database-demo/services"
)

func TestGetShardID(t *testing.T) {
    services.InitDatabase()
    id := services.GetShardID("user1")
    if id < 0 || id >= 3 {
        t.Errorf("Invalid shard ID: %d", id)
    }
}