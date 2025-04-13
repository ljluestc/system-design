package storage

import (
    "github.com/go-redis/redis/v8"
    "context"
    "crawler/internal/utils"
)

type Redis struct {
    client *redis.Client
}

func NewRedis(addr string) (*Redis, error) {
    client := redis.NewClient(&redis.Options{Addr: addr})
    _, err := client.Ping(context.Background()).Result()
    if err != nil {
        return nil, utils.Errorf("connect redis: %v", err)
    }
    return &Redis{client: client}, nil
}

func (s *Redis) Set(key string, value []byte) error {
    return s.client.Set(context.Background(), key, value, 0).Err()
}

func (s *Redis) Get(key string) ([]byte, error) {
    val, err := s.client.Get(context.Background(), key).Bytes()
    if err == redis.Nil {
        return nil, utils.Errorf("key %s not found", key)
    }
    return val, err
}