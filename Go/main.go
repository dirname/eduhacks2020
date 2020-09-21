package main

import (
	"eduhacks2020/Go/define"
	"eduhacks2020/Go/middleware"
	"eduhacks2020/Go/pkg/etcd"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/protocol/websocket"
	"eduhacks2020/Go/routers"
	"eduhacks2020/Go/utils"
	"eduhacks2020/Go/utils/log"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
	"net"
)

func init() {
	setting.ReadConfigure()
	log.Setup()
}

func initRPCServer() {
	//如果是集群，则启用RPC进行通讯
	if utils.IsCluster() {
		//初始化RPC服务
		websocket.InitGRPCServer()
		fmt.Printf("Start RPC Listening on :%s\n", setting.CommonSetting.RPCPort)
	}
}

//ETCD注册发现服务
func registerServer() {
	if utils.IsCluster() {
		// 注册租约
		ser, err := etcd.NewServiceReg(setting.EtcdSetting.Endpoints, 5)
		if err != nil {
			panic(err)
		}

		hostPort := net.JoinHostPort(setting.GlobalSetting.LocalHost, setting.CommonSetting.RPCPort)
		// 添加key
		err = ser.PutService(define.EtcdServerList+hostPort, hostPort)
		if err != nil {
			panic(err)
		}

		cli, err := etcd.NewClientDis(setting.EtcdSetting.Endpoints)
		if err != nil {
			panic(err)
		}
		_, err = cli.GetService(define.EtcdServerList)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	gin.SetMode(gin.DebugMode) // 生产模式中改写成 release

	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.CSRF())
	router.Use(middleware.Auth()) // 授权的中间件
	//  RPC 初始
	initRPCServer()
	// 注册 etcd
	registerServer()
	// 初始化路由
	dm := routers.Init(router)
	defer dm.Close()
	// 定时发送心跳包
	websocket.PingTimer()
	//database.DefaultSetting()
	// 数据库的一些初始化, 全局引用

	//if err := redis.Instance.Set(context.Background(), "AdminUser", setting.AdminConf.Username, 0).Err(); err != nil {
	//	log.Errorf("Set Admin User Failed: %s", err.Error())
	//}
	//
	//if err := redis.Instance.Set(context.Background(), "AdminPassword", setting.AdminConf.Password, 0).Err(); err != nil {
	//	log.Errorf("Set Admin Password Failed: %s", err.Error())
	//}

	//defer redis.Close()

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
