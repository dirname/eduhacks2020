package websocket

import (
	"eduhacks2020/Go/api/users"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	websocket2 "eduhacks2020/Go/routers/websocket"
	"eduhacks2020/Go/utils"
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
	UserName    string          // 业务端标识用户账号
	NickName    string          // 业务端标识用户昵称
	UserRole    int             // 业务端标识用户角色
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
				if c.NickName == "" || c.UserName == "" || c.UserRole == 0 {
					c.setInfo(id)
				}
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

// setInfo 设置业务端信息
func (c *Client) setInfo(sessionID string) {
	res := &protobuf.Response{
		Code:   -1,
		Msg:    users.TokenInvalid,
		Type:   3,
		Data:   nil,
		Render: false,
		Html:   nil,
		Id:     "",
	}
	logout, _ := proto.Marshal(res)
	session := database.SessionManager{Values: make(map[interface{}]interface{})}
	data, err := session.GetData(sessionID)
	if err != nil {
		return
	}
	session.DecryptedData(data.(string), database.SessionName)
	token := session.Values["token"]
	if token == nil {
		c.Socket.WriteMessage(2, xorData(logout, false))
		return
	}
	claims, err := utils.ParseToken(token.(string))
	if err != nil {
		c.Socket.WriteMessage(2, xorData(logout, false))
		return
	}
	c.UserName = claims.Username
	c.NickName = claims.Name
	c.UserRole = claims.Role
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
	router := websocket2.Router{}
	param := &websocket2.ProtoParam{
		Request:   &req,
		Response:  res,
		SessionID: id,
		Redis:     r.Instance,
		DB:        d.DB,
	}
	router.Find(param, websocket2.Handler)
	data, _ := proto.Marshal(res)
	c.Socket.WriteMessage(2, xorData(data, false))
}
