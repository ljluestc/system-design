// node/replica.go
package node

import (
    "key-value-store/logger"
)

func (n *Node) DemoteToReplica(primary *Node) bool {
    if !n.IsPrimary {
        return false
    }
    n.mutex.Lock()
    defer n.mutex.Unlock()
    n.IsPrimary = false
    n.Replicator = nil
    n.Primary = primary
    logger.Info("Node demoted to replica", "shard", n.ShardID)
    return true
}