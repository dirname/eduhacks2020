package main

import (
	"context"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/middleware"
	"eduhacks2020/Go/protocol/websocket"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func main() {
	gin.SetMode(gin.DebugMode) // 生产模式中改写成 release

	websocketHandler := &protocol.Controller{}
	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.CSRF())
	router.Use(middleware.Auth()) // 授权的中间件

	//database.ReadConfigure() // 读取配置
	database.DefaultSetting()
	// 数据库的一些初始化, 全局引用
	orm := database.ORM{}
	orm.Init()
	defer orm.Close()

	redis := database.RedisClient{}
	redis.Init()

	if err := redis.Instance.Set(context.Background(), "AdminUser", database.AdminConf.Username, 0).Err(); err != nil {
		log.Errorf("Set Admin User Failed: %s", err.Error())
	}

	if err := redis.Instance.Set(context.Background(), "AdminPassword", database.AdminConf.Password, 0).Err(); err != nil {
		log.Errorf("Set Admin Password Failed: %s", err.Error())
	}

	defer redis.Close()

	router.GET("/ws", func(c *gin.Context) {
		websocketHandler.Run(c.Writer, c.Request, &orm, &redis)
	})

	//router.POST("/sessions", func(c *gin.Context) {
	//	// 创建一个 sessions 存储容器
	//	store, dbSession := database.CreateMongoStore()
	//	defer dbSession.Close()
	//	session, err := store.Get(c.Request, database.SessionName)
	//	if err != nil {
	//		log.Error(err.Error())
	//	}
	//	c.JSON(200, gin.H{"count": session.Values["pv"], "id": session.ID})
	//})
	//router.GET("/login", func(context *gin.Context) {
	//	store, dbSession := database.CreateMongoStore()
	//	defer dbSession.Close()
	//	session, err := store.Get(context.Request, database.SessionName)
	//	if err != nil {
	//		log.Error(err.Error())
	//	}
	//	session.Values["login"] = true
	//	session.Save(context.Request, context.Writer)
	//	context.JSON(200, gin.H{"msg": "ok", "id": session.ID})
	//})
	//router.GET("/logout", func(context *gin.Context) {
	//	store, dbSession := database.CreateMongoStore()
	//	defer dbSession.Close()
	//	session, err := store.Get(context.Request, database.SessionName)
	//	if err != nil {
	//		log.Error(err.Error())
	//	}
	//	session.Values["login"] = true
	//	session.Save(context.Request, context.Writer)
	//	context.JSON(200, gin.H{"msg": "ok", "id": session.ID})
	//})
	router.Run(":555")
}
