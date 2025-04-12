package cache

type ReplicationManager struct {
    hashRing *HashRing
    nodes    map[string]*CacheNode
}

func NewReplicationManager(hashRing *HashRing, nodes map[string]*CacheNode) *ReplicationManager {
    return &ReplicationManager{
        hashRing: hashRing,
        nodes:    nodes,
    }
}

func (r *ReplicationManager) ReplicateSet(key, value string, ttl time.Duration) {
    nodeAddr := r.hashRing.GetNode(key)
    node := r.nodes[nodeAddr]
    node.Set(key, value, ttl)
}