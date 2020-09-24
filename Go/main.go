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
	"os"
	"os/signal"
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
	router.Use(middleware.CORS()) // 跨域
	router.Use(middleware.Logger())
	router.Use(middleware.CSRF())
	router.Use(middleware.Auth()) // 授权的中间件
	//  RPC 初始
	initRPCServer()
	// 注册 etcd
	registerServer()
	// 初始化路由
	dm := routers.DatabaseManager{}
	dm.Init(router)
	// 捕获关闭的信号
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		dm.Close()
		os.Exit(1)
	}()
	// 定时发送心跳包
	websocket.PingTimer()

	router.Run(":" + setting.CommonSetting.Port)
}
