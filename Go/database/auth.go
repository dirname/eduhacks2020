package database

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// RedisAuth 使用 Redis 来鉴权
type RedisAuth struct {
	Instance *redis.Client
}

var ctx = context.Background()

// Init 初始化 Redis 的连接
func (r *RedisAuth) Init() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPwd, // no password set
		DB:       0,        // use default DB
	})
	r.Instance = rdb
}

// Close 关闭 Redis 的连接
func (r *RedisAuth) Close() {
	r.Instance.Close()
}

// GetFlag 获取标识
func (r *RedisAuth) GetFlag(id string) (string, error) {
	return r.Instance.Get(ctx, id).Result()
}

// SetFlag 设置标识
func (r *RedisAuth) SetFlag(id, flag string) error {
	return r.Instance.Set(ctx, id, flag, 0).Err()
}
