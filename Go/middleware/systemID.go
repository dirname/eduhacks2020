package middleware

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/define"
	"eduhacks2020/Go/define/retcode"
	"eduhacks2020/Go/pkg/etcd"
	"eduhacks2020/Go/protocol/websocket"
	"eduhacks2020/Go/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SystemIDMiddleware 的中间件
func SystemIDMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.Method != http.MethodPost {
			context.Writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		//检查header是否设置SystemId
		systemID := context.Request.Header.Get("SystemID")
		if len(systemID) == 0 {
			context.Abort()
			api.Render(context.Writer, retcode.FAIL, "SystemID cannot be empty", []string{})
			return
		}

		//判断是否被注册
		if utils.IsCluster() {
			resp, err := etcd.Get(define.EtcdPrefixAccountInfo + systemID)
			if err != nil {
				context.Abort()
				api.Render(context.Writer, retcode.FAIL, "Etcd server error", []string{})
				return
			}

			if resp.Count == 0 {
				context.Abort()
				api.Render(context.Writer, retcode.FAIL, "Invalid SystemID", []string{})
				return
			}
		} else {
			if _, ok := websocket.SystemMap.Load(systemID); !ok {
				context.Abort()
				api.Render(context.Writer, retcode.FAIL, "Invalid SystemID", []string{})
				return
			}
		}
		context.Next()
	}
}
