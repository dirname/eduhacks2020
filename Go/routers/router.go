package routers

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/protocol/websocket"
	"github.com/gin-gonic/gin"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	ORM   *database.ORM
	Redis *database.RedisClient
}

// Init 初始化的路由
func (d *DatabaseManager) Init(engine *gin.Engine) {

	orm := database.ORM{}
	orm.Init()

	redis := database.RedisClient{}
	redis.Init()

	d.ORM = &orm
	d.Redis = &redis

	engine.GET("/get/client", func(context *gin.Context) {
		context.JSON(200, websocket.Manager.AllClient())
	})

	websocket.StartWebSocket(engine, &orm, &redis)
	go websocket.WriteMessage()
}

// Close 关闭数据库连接
func (d *DatabaseManager) Close() {
	d.ORM.Close()
	d.Redis.Close()
}
