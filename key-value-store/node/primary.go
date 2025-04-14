// node/primary.go
package node

import (
    "key-value-store/logger"
)

func (n *Node) PromoteToPrimary() bool {
    if n.IsPrimary {
        return false
    }
    n.mutex.Lock()
    defer n.mutex.Unlock()
    n.IsPrimary = true
    n.Replicator = replication.NewReplicator(n)
    n.Primary = n
    logger.Info("Node promoted to primary", "shard", n.ShardID)
    return true
}