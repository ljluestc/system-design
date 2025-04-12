package db

import (
    "key-value-store/internal/config"
    "key-value-store/internal/utils"
)

type Node struct {
    ID       string
    Capacity int
    Storage  Storage
}

var nodes []*Node

func InitNodes(cfg *config.Config) error {
    for i := 0; i < cfg.NodeCount; i++ {
        nodes = append(nodes, &Node{
            ID:       utils.GenerateUUID(),
            Capacity: cfg.DefaultCapacity,
            Storage:  NewMemoryStorage(),
        })
    }
    return nil
}

func Get(key string) (string, error) {
    node := GetResponsibleNode(key)
    return node.Storage.Get(key)
}

func Put(key, value string) error {
    node := GetResponsibleNode(key)
    return node.Storage.Put(key, value)
}