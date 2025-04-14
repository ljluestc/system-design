package inference

import (
    "chatgpt/config"
    "context"
)

type Service struct {
    shards []*Shard
}

func NewService(numShards int) *Service {
    s := &Service{shards: make([]*Shard, numShards)}
    for i := 0; i < numShards; i++ {
        s.shards[i] = NewShard(i)
    }
    return s
}

func (s *Service) ProcessQuery(userID, prompt string) (string, error) {
    shardID := GetShardID(userID, config.GetConfig().NumShards)
    shard := s.shards[shardID]
    ctx, _ := context.Background(), "" // Simplified; real context would be fetched
    response := shard.Model.Generate(prompt, ctx)
    return response, nil
}