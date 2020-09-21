package database

import (
	"context"
	"eduhacks2020/Go/pkg/setting"
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
		Addr:     setting.Database.RedisHost,
		Password: setting.Database.RedisPwd, // no password set
		DB:       0,                         // use default DB
	})
	r.Instance = rdb
}

// Close 关闭 Redis 的连接
func (r *RedisClient) Close() {
	r.Instance.Close()
}
