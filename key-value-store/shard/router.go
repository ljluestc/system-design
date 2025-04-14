// shard/router.go
package shard

import (
    "key-value-store/node"
    "key-value-store/utils"
)

type Router struct {
    nodes []*node.Node
}

func NewRouter(nodes []*node.Node) *Router {
    return &Router{nodes: nodes}
}

func (r *Router) Route(key string) *node.Node {
    shardID := GetShardID(key)
    for _, n := range r.nodes {
        if n.ShardID == shardID && n.IsPrimary {
            return n
        }
    }
    utils.Error("No primary found", "shard", shardID)
    return nil
}