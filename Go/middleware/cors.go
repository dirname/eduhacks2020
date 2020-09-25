package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS() gin.HandlerFunc {
	return func(context *gin.Context) {
		context.Header("Access-Control-Allow-Origin", context.Request.Header.Get("Origin"))
		context.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token, SystemID")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		//context.Header("Access-Control-Allow-Credentials", "true") 使用 Header 认证不需要 session, session 为了调试方便使用

		if context.Request.Method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}

		context.Next()
	}
}
