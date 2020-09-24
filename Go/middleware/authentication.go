package middleware

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/models"
	"eduhacks2020/Go/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Auth 认证是否登陆的中间件
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		store, dbSession := database.CreateMongoStore()
		defer dbSession.Close()
		session, err := store.Get(context.Request, database.SessionName)
		if err != nil {
			log.Error(err)
		}

		//这里为了避免 session 没有存到 Mongo 中, 所以每次刷新页面, 将会保存一次确保在表里
		var count int
		pv := session.Values["pv"]
		if pv == nil {
			count = 1
		} else {
			count = pv.(int)
			count++
		}
		token := session.Values["token"]
		session.Values["pv"] = count
		session.Save(context.Request, context.Writer)

		path := context.Request.URL.Path
		clientIP := context.ClientIP()
		method := context.Request.Method
		tokenStr := ""
		isLogin := session.Values["login"]
		if path != "/ws" { //如果是 websocket 的连接, 那么将在 websocket 里开始认证
			res := models.Response{
				Code:   http.StatusUnauthorized,
				Path:   path,
				Method: method,
				Msg:    "Unauthorized",
				Data:   nil,
				Time:   time.Now(),
				IP:     clientIP,
			}
			if token == nil {
				tokenStr = context.Request.Header.Get("Authorization")
				if tokenStr != "" {
					isLogin = true
				} else {
					context.Abort()
					context.JSON(http.StatusUnauthorized, res)
					return
				}
			} else {
				tokenStr = token.(string)
			}
			if isLogin == false || isLogin == nil {
				context.Abort()
				context.JSON(http.StatusUnauthorized, res)
				return
			}
			login, role := Validate(tokenStr)
			if !login {
				context.Abort()
				context.JSON(http.StatusUnauthorized, res)
				return
			}
			if !role {
				res.Code = http.StatusForbidden
				res.Method = "You do not have permission to view the page"
				context.Abort()
				context.JSON(http.StatusForbidden, res)
				return
			}
		}
		context.Next()
	}
}

// Validate 认证一些用户的信息
func Validate(tokenString string) (bool, bool) {
	claims, err := utils.ParseToken(tokenString)
	if err != nil {
		return false, false
	}
	redis := database.RedisClient{}
	redis.Init()
	defer redis.Close()
	redisAuth := api.AuthRedis{Redis: redis.Instance}
	flag, _ := redisAuth.GetFlag(claims.UID)
	if claims.Flag != flag {
		return false, false
	}
	if claims.Role != -1 {
		return true, false
	}
	return true, true
}
