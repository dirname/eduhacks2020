package websocket

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"time"
)

// ToClientChan channel通道
var ToClientChan chan clientInfo

//channel通道结构体
type clientInfo struct {
	ClientID   string
	SendUserID string
	MessageID  string
	Code       int
	Msg        string
	Data       *string
}

// RetData
type RetData struct {
	MessageID  string      `json:"messageId"`
	SendUserID string      `json:"sendUserId"`
	Code       int         `json:"code"`
	Msg        string      `json:"msg"`
	Data       interface{} `json:"data"`
}

// 心跳间隔
var heartbeatInterval = 25 * time.Second

func init() {
	ToClientChan = make(chan clientInfo, 1000)
}

// Manager
var Manager = NewClientManager() // 管理者

// StartWebSocket
func StartWebSocket(engine *gin.Engine, o *database.ORM, r *database.RedisClient) {
	websocketHandler := &Controller{}
	engine.GET("/ws", func(c *gin.Context) {
		websocketHandler.Run(c.Writer, c.Request, o, r)
	})
	//websocketHandler := &Controller{}
	//http.HandleFunc("/ws", websocketHandler.Run)
	go Manager.Start()
}

// SendMessage2Client 发送信息到指定客户端
func SendMessage2Client(clientID string, sendUserId string, code int, msg string, data *string) (messageId string) {
	messageId = utils.GenUUID()
	if utils.IsCluster() {
		addr, _, _, isLocal, err := utils.GetAddrInfoAndIsLocal(clientID)
		if err != nil {
			log.Errorf("%s", err)
			return
		}

		//如果是本机则发送到本机
		if isLocal {
			SendMessage2LocalClient(messageId, clientID, sendUserId, code, msg, data)
		} else {
			//发送到指定机器
			SendRPC2Client(addr, messageId, sendUserId, clientID, code, msg, data)
		}
	} else {
		//如果是单机服务，则只发送到本机
		SendMessage2LocalClient(messageId, clientID, sendUserId, code, msg, data)
	}

	return
}

// CloseClient 关闭客户端
func CloseClient(clientID, systemId string) {
	if utils.IsCluster() {
		addr, _, _, isLocal, err := utils.GetAddrInfoAndIsLocal(clientID)
		if err != nil {
			log.Errorf("%s", err)
			return
		}

		//如果是本机则发送到本机
		if isLocal {
			CloseLocalClient(clientID, systemId)
		} else {
			//发送到指定机器
			CloseRPCClient(addr, clientID, systemId)
		}
	} else {
		//如果是单机服务，则只发送到本机
		CloseLocalClient(clientID, systemId)
	}

	return
}

// AddClient2Group 添加客户端到分组
func AddClient2Group(systemID string, groupName string, clientId string, userId string, extend string) {
	//如果是集群则用redis共享数据
	if utils.IsCluster() {
		//判断key是否存在
		addr, _, _, isLocal, err := utils.GetAddrInfoAndIsLocal(clientId)
		if err != nil {
			log.Errorf("%s", err)
			return
		}

		if isLocal {
			if client, err := Manager.GetByClientID(clientId); err == nil {
				//添加到本地
				Manager.AddClient2LocalGroup(groupName, client, userId, extend)
			} else {
				log.Error(err)
			}
		} else {
			//发送到指定的机器
			SendRPCBindGroup(addr, systemID, groupName, clientId, userId, extend)
		}
	} else {
		if client, err := Manager.GetByClientID(clientId); err == nil {
			//如果是单机，就直接添加到本地group了
			Manager.AddClient2LocalGroup(groupName, client, userId, extend)
		}
	}
}

// SendMessage2Group 发送信息到指定分组
func SendMessage2Group(systemID, sendUserId, groupName string, code int, msg string, data *string) (messageID string) {
	messageID = utils.GenUUID()
	if utils.IsCluster() {
		//发送分组消息给指定广播
		go SendGroupBroadcast(systemID, messageID, sendUserId, groupName, code, msg, data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalGroup(systemID, messageID, sendUserId, groupName, code, msg, data)
	}
	return
}

// SendMessage2System 发送信息到指定系统
func SendMessage2System(systemID, sendUserId string, code int, msg string, data string) {
	messageID := utils.GenUUID()
	if utils.IsCluster() {
		//发送到系统广播
		SendSystemBroadcast(systemID, messageID, sendUserId, code, msg, &data)
	} else {
		//如果是单机服务，则只发送到本机
		Manager.SendMessage2LocalSystem(systemID, messageID, sendUserId, code, msg, &data)
	}
}

// GetOnlineList 获取分组列表
func GetOnlineList(systemID *string, groupName *string) map[string]interface{} {
	var clientList []string
	if utils.IsCluster() {
		//发送到系统广播
		clientList = GetOnlineListBroadcast(systemID, groupName)
	} else {
		//如果是单机服务，则只发送到本机
		retList := Manager.GetGroupClientList(utils.GenGroupKey(*systemID, *groupName))
		clientList = append(clientList, retList...)
	}

	return map[string]interface{}{
		"count": len(clientList),
		"list":  clientList,
	}
}

// SendMessage2LocalClient 通过本服务器发送信息
func SendMessage2LocalClient(messageID, clientId string, sendUserId string, code int, msg string, data *string) {
	log.WithFields(log.Fields{
		"host":     setting.GlobalSetting.LocalHost,
		"port":     setting.CommonSetting.Port,
		"clientId": clientId,
	}).Info("发送到通道")
	ToClientChan <- clientInfo{ClientID: clientId, MessageID: messageID, SendUserID: sendUserId, Code: code, Msg: msg, Data: data}
	return
}

// CloseLocalClient 发送关闭信号
func CloseLocalClient(clientID, systemId string) {
	if conn, err := Manager.GetByClientID(clientID); err == nil && conn != nil {
		if conn.SystemID != systemId {
			return
		}
		Manager.DisConnect <- conn
		log.WithFields(log.Fields{
			"host":     setting.GlobalSetting.LocalHost,
			"port":     setting.CommonSetting.Port,
			"clientID": clientID,
		}).Info("主动踢掉客户端")
	}
	return
}

// WriteMessage 监听并发送给客户端信息
func WriteMessage() {
	for {
		clientInfo := <-ToClientChan
		log.WithFields(log.Fields{
			"host":       setting.GlobalSetting.LocalHost,
			"port":       setting.CommonSetting.Port,
			"clientId":   clientInfo.ClientID,
			"messageId":  clientInfo.MessageID,
			"sendUserId": clientInfo.SendUserID,
			"code":       clientInfo.Code,
			"msg":        clientInfo.Msg,
			"data":       clientInfo.Data,
		}).Info("发送到本机")
		if conn, err := Manager.GetByClientID(clientInfo.ClientID); err == nil && conn != nil {
			if err := Render(conn.Socket, clientInfo.MessageID, clientInfo.SendUserID, clientInfo.Code, clientInfo.Msg, clientInfo.Data); err != nil {
				Manager.DisConnect <- conn
				log.WithFields(log.Fields{
					"host":     setting.GlobalSetting.LocalHost,
					"port":     setting.CommonSetting.Port,
					"clientId": clientInfo.ClientID,
					"msg":      clientInfo.Msg,
				}).Error("客户端异常离线：" + err.Error())
			}
		}
	}
}

// Render 渲染
func Render(conn *websocket.Conn, messageID string, sendUserId string, code int, message string, data interface{}) error {
	return conn.WriteJSON(RetData{
		Code:       code,
		MessageID:  messageID,
		SendUserID: sendUserId,
		Msg:        message,
		Data:       data,
	})
}

// PingTimer 启动定时器进行心跳检测
func PingTimer() {
	go func() {
		ticker := time.NewTicker(heartbeatInterval)
		defer ticker.Stop()
		for {
			<-ticker.C
			//发送心跳
			for clientID, conn := range Manager.AllClient() {
				log.Info("心跳包来了")
				if err := conn.Socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(time.Second)); err != nil {
					Manager.DisConnect <- conn
					log.Errorf("发送心跳失败: %s 总连接数：%d", clientID, Manager.Count())
				}
			}
		}

	}()
}
