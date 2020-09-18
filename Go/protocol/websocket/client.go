package protocol

import (
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Client websocket 的客户端结构
type Client struct {
	ClientID    string          // 标识ID
	SystemID    string          // 系统ID
	Socket      *websocket.Conn // 用户连接
	ConnectTime uint64          // 首次连接时间
	IsDeleted   bool            // 是否删除或下线
	UserID      string          // 业务端标识用户ID
	Extend      string          // 扩展字段，用户可以自定义
	GroupList   []string
}

// SendData 发送消息的结构体
type SendData struct {
	Code int
	Msg  string
	Data *interface{}
}

// NewClient 创建一个新的 websocket 客户端
func NewClient(clientID string, systemID string, socket *websocket.Conn) *Client {
	return &Client{
		ClientID:    clientID,
		SystemID:    systemID,
		Socket:      socket,
		ConnectTime: uint64(time.Now().Unix()),
		IsDeleted:   false,
	}
}

// Read 客户端读取消息
func (c *Client) Read(d *database.ORM, r2 *database.RedisClient, id string) {
	go func() {
		for {
			messageType, msg, err := c.Socket.ReadMessage()
			if err != nil {
				if messageType == -1 && websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure, websocket.CloseNoStatusReceived) {
					//Manager.DisConnect <- c
					return
				} else if messageType != websocket.PingMessage {
					return
				}
			} else {
				c.Router(msg, d, r2, id)
			}
		}
	}()
}

func xorData(data []byte, decrypt bool) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = b ^ 32
	}
	if !decrypt {
		return []byte(base64.URLEncoding.EncodeToString(res))
	}
	return res
}

// Router 客户端处理路由
func (c *Client) Router(msg []byte, d *database.ORM, r *database.RedisClient, id string) {
	msg = xorData(msg, true)
	var req protobuf.Request
	err := proto.Unmarshal(msg, &req)
	res := &protobuf.Response{
		Code:   http.StatusInternalServerError,
		Msg:    "Invalid requests",
		Type:   3,
		Data:   nil,
		Render: true,
		Html: &protobuf.Render{
			Code:   render.GetLayer(0, render.Incorrect, "Error", "Invalid Requests"),
			Type:   1,
			Id:     "layerMsgBox",
			Iframe: false,
		},
		Id: "",
	}
	if err != nil {
		log.Error("Parse Protobuf Error: ", err.Error(), string(msg))
		data, _ := proto.Marshal(res)
		c.Socket.WriteMessage(2, xorData(data, false))
		return
	}
	router := database.Router{}
	router.Find(&database.ProtoParam{
		Request:   &req,
		Response:  res,
		SessionID: id,
		Redis:     r.Instance,
		DB:        d.DB,
	}, database.Handler)
	data, _ := proto.Marshal(res)
	c.Socket.WriteMessage(2, xorData(data, false))
}
