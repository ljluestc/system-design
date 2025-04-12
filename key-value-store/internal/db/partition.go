package db

import "key-value-store/internal/utils"

func GetResponsibleNode(key string) *Node {
    hash := utils.HashKey(key)
    nodeIndex := hash % uint32(len(nodes))
    return nodes[nodeIndex]
}