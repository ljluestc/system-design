package services

import "github.com/go-redis/redis/v8"

var redisClient = redis.NewClient(&redis.Options{Addr: "localhost:6379"})

func CacheSet(key, value string) error {
    return redisClient.Set(ctx, key, value, 0).Err()
}

func CacheGet(key string) (string, error) {
    return redisClient.Get(ctx, key).Result()
}