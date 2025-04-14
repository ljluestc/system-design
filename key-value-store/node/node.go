// node/node.go
package node

import (
    "key-value-store/logger"
    "key-value-store/replication"
    "key-value-store/shard"
    "key-value-store/storage"
    "sync"
)

type Node struct {
    ShardID    int
    IsPrimary  bool
    Storage    storage.Storage
    Replicas   []*Node
    Replicator *replication.Replicator
    Primary    *Node
    mutex      sync.RWMutex
}

func NewNode(shardID int, isPrimary bool, primary *Node) *Node {
    n := &Node{
        ShardID:   shardID,
        IsPrimary: isPrimary,
        Storage:   storage.NewMemoryStorage(),
        Replicas:  make([]*Node, 0),
    }
    if isPrimary {
        n.Replicator = replication.NewReplicator(n)
    } else {
        n.Primary = primary
    }
    logger.Info("Node created", "shard", shardID, "primary", isPrimary)
    return n
}

func (n *Node) Put(key string, value interface{}) bool {
    if shard.GetShardID(key) != n.ShardID {
        logger.Warn("Wrong shard", "key", key, "shard", n.ShardID)
        return false
    }
    if !n.IsPrimary {
        return n.Primary.Put(key, value)
    }
    n.mutex.Lock()
    defer n.mutex.Unlock()
    n.Storage.Put(key, value)
    n.Replicator.Replicate(key, value)
    logger.Info("Stored", "key", key, "shard", n.ShardID)
    return true
}

func (n *Node) Get(key string) (interface{}, bool) {
    if shard.GetShardID(key) != n.ShardID {
        return nil, false
    }
    n.mutex.RLock()
    defer n.mutex.RUnlock()
    return n.Storage.Get(key)
}

func (n *Node) Delete(key string) bool {
    if shard.GetShardID(key) != n.ShardID {
        return false
    }
    if !n.IsPrimary {
        return n.Primary.Delete(key)
    }
    n.mutex.Lock()
    defer n.mutex.Unlock()
    ok := n.Storage.Delete(key)
    if ok {
        n.Replicator.ReplicateDelete(key)
    }
    return ok
}

func (n *Node) AddReplica(replica *Node) {
    if n.IsPrimary {
        n.Replicas = append(n.Replicas, replica)
        n.Replicator.AddReplica(replica)
    }
}