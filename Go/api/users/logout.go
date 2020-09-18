package users

import (
	"context"
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/utils"
	"errors"
	"github.com/go-redis/redis/v8"
)

type LogoutParam struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

func (l *LogoutParam) Exec(redis *redis.Client) ([]byte, string, error) {
	if claims, err := utils.ParseToken(l.Token); err != nil {
		return nil, TokenInvalid, errors.New(TokenInvalid)
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			return nil, TokenInvalid, errors.New(TokenInvalid)
		}
		redis.Del(context.Background(), claims.UID).Result()
		return nil, "Logout success !", nil
	}
}
