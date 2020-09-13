package database

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// RedisClient 使用 Redis 来鉴权
type RedisClient struct {
	Instance *redis.Client
}

var ctx = context.Background()

// Init 初始化 Redis 的连接
func (r *RedisClient) Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPwd, // no password set
		DB:       0,        // use default DB
	})
	r.Instance = rdb
}

// Close 关闭 Redis 的连接
func (r *RedisClient) Close() {
	r.Instance.Close()
}
