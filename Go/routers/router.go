package routers

import (
	"eduhacks2020/Go/api/bind2group"
	"eduhacks2020/Go/api/closeclient"
	"eduhacks2020/Go/api/control"
	"eduhacks2020/Go/api/getonlinelist"
	"eduhacks2020/Go/api/register"
	"eduhacks2020/Go/api/send2client"
	"eduhacks2020/Go/api/send2clients"
	"eduhacks2020/Go/api/send2group"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/middleware"
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
	registerHandler := &register.Controller{}
	sendToClientHandler := &send2client.Controller{}
	sendToClientsHandler := &send2clients.Controller{}
	sendToGroupHandler := &send2group.Controller{}
	bindToGroupHandler := &bind2group.Controller{}
	getGroupListHandler := &getonlinelist.Controller{}
	closeClientHandler := &closeclient.Controller{}

	engine.GET(APIManagerClientGet, c.GetDevices)
	engine.POST(APIManagerRegister, registerHandler.Run)
	engine.POST(APIManagerMsgSend, middleware.SystemIDMiddleware(), sendToClientHandler.Run)
	engine.POST(APIManagerCloseClient, middleware.SystemIDMiddleware(), closeClientHandler.Run)
	engine.POST(APIManagerMsgSends, middleware.SystemIDMiddleware(), sendToClientsHandler.Run)
	engine.POST(APIManagerGroupSend, middleware.SystemIDMiddleware(), sendToGroupHandler.Run)
	engine.POST(APIManagerGroupBind, middleware.SystemIDMiddleware(), bindToGroupHandler.Run)
	engine.POST(APIManagerGroupGet, middleware.SystemIDMiddleware(), getGroupListHandler.Run)

	websocket.StartWebSocket(engine, &orm, &redis, &mongo, d.Mongo.CollectionName)
	go websocket.WriteMessage()
}

// Close 关闭数据库连接
func (d *DatabaseManager) Close() {
	d.ORM.Close()
	d.Redis.Close()
	d.Mongo.Close() // 删除客户端连接表
}
