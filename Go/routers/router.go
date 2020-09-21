package routers

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/protocol/websocket"
	"github.com/gin-gonic/gin"
)

type databaseManager struct {
	ORM   *database.ORM
	Redis *database.RedisClient
}

// Init 初始化的路由
func Init(engine *gin.Engine) *databaseManager {

	orm := database.ORM{}
	orm.Init()

	redis := database.RedisClient{}
	redis.Init()

	dm := databaseManager{
		ORM:   &orm,
		Redis: &redis,
	}

	engine.GET("/get/client", func(context *gin.Context) {
		context.JSON(200, websocket.Manager.AllClient())
	})

	websocket.StartWebSocket(engine, &orm, &redis)
	go websocket.WriteMessage()
	return &dm
}

// Close
func (d *databaseManager) Close() {
	d.ORM.Close()
	d.Redis.Close()
}
