package routers

import (
	"eduhacks2020/Go/api/control"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/protocol/websocket"
	"github.com/gin-gonic/gin"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	ORM   *database.ORM
	Redis *database.RedisClient
	Mongo *database.MongoClientDevice
}

// Init 初始化的路由
func (d *DatabaseManager) Init(engine *gin.Engine) {

	orm := database.ORM{}
	orm.Init()

	redis := database.RedisClient{}
	redis.Init()

	mongo := database.MongoClientDevice{}
	mongo.Init()

	d.ORM = &orm
	d.Redis = &redis
	d.Mongo = &mongo

	c := control.DevicesControl{Mongo: d.Mongo}

	engine.GET(APIManagerClientGet, c.GetDevices)

	websocket.StartWebSocket(engine, &orm, &redis, &mongo, d.Mongo.CollectionName)
	go websocket.WriteMessage()
}

// Close 关闭数据库连接
func (d *DatabaseManager) Close() {
	d.ORM.Close()
	d.Redis.Close()
	d.Mongo.Close() // 删除客户端连接表
}
