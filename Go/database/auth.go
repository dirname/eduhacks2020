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

func (r *RedisAuth) GetFlag(id string) (string, error) {
	return r.Instance.Get(ctx, id).Result()
}
