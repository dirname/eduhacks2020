package main

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/middleware"
	"eduhacks2020/Go/protocol/websocket"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
)

func main() {
	gin.SetMode(gin.DebugMode) //生产模式中改写成 release

	websocketHandler := &protocol.Controller{}
	router := gin.Default()

	router.Use(middleware.Logger())
	router.Use(middleware.CSRF())
	router.Use(middleware.Auth()) //授权的中间件
	orm := database.ORM{}
	orm.Init()
	defer orm.Close()
	router.GET("/ws", func(c *gin.Context) {
		websocketHandler.Run(c.Writer, c.Request, &orm)
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
	//http.HandleFunc("/ws", websocketHandler.Run)
	//if err := http.ListenAndServe(":555", nil); err != nil {
	//	panic(err)
	//}
}

//func helloWorld(w http.ResponseWriter, r *http.Request) {
//	test := &protobuf.Student{
//		Name:   "Hao",
//		Male:   true,
//		Scores: []int32{98, 85, 88},
//	}
//	data, err := proto.Marshal(test)
//	if err != nil {
//		return
//	}
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Set("content-type", "application/octet-stream")
//	w.Write(data)
//}
