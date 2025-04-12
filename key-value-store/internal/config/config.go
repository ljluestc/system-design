package config

import "errors"

type Config struct {
    Port             string
    NodeCount        int
    DefaultCapacity  int
    ReplicationFactor int
    ConsistencyLevel string
}

func LoadConfig() (*Config, error) {
    return &Config{
        Port:             "8080",
        NodeCount:        10,
        DefaultCapacity:  100,
        ReplicationFactor: 3,
        ConsistencyLevel: "eventual",
    }, nil
}