package api

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// AuthRedis 用于认证的 redis 实例
type AuthRedis struct {
	Redis *redis.Client
}

var ctx = context.Background()

// GetFlag 获取标识
func (r *AuthRedis) GetFlag(id string) (string, error) {
	return r.Redis.Get(ctx, id).Result()
}

// SetFlag 设置标识
func (r *AuthRedis) SetFlag(id, flag string) error {
	return r.Redis.Set(ctx, id, flag, 0).Err()
}