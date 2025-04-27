package cache

import (
    "context"
    "time"

    "github.com/go-redis/redis/v8"
)

type RedisClient = redis.Client

func NewRedis(uri string) *RedisClient {
    opt, err := redis.ParseURL(uri)
    if err != nil {
        panic(err)
    }
    rdb := redis.NewClient(opt)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := rdb.Ping(ctx).Err(); err != nil {
        panic(err)
    }
    return rdb
}
