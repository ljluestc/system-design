package inference

import (
    "crypto/md5"
    "strconv"
)

type Shard struct {
    Model *Model
}

func NewShard(id int) *Shard {
    return &Shard{Model: NewModel(id)}
}

func GetShardID(key string, numShards int) int {
    hash := md5.Sum([]byte(key))
    hashInt, _ := strconv.ParseInt(string(hash[:8]), 16, 64)
    return int(hashInt) % numShards
}