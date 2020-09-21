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

// 初始化的路由
func Init(engine *gin.Engine) *databaseManager {

	orm := database.ORM{}
	orm.Init()

	redis := database.RedisClient{}
	redis.Init()

	dm := databaseManager{
		ORM:   &orm,
		Redis: &redis,
	}

	websocket.StartWebSocket(engine, &orm, &redis)
	go websocket.WriteMessage()
	return &dm
}

func (d *databaseManager) Close() {
	d.ORM.Close()
	d.Redis.Close()
}
