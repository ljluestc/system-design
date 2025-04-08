package unit

import (
    "testing"
    "database-demo/services"
)

func TestSyncReplicas(t *testing.T) {
    if err := services.SyncReplicas(); err != nil {
        t.Errorf("Sync failed: %v", err)
    }
}