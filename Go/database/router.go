package database

import (
	"eduhacks2020/Go/models/request"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	"encoding/json"
	"gorm.io/gorm"
	"net/http"
)

// 定义路由
const (
	APILogin = "/api/login" //用户登录的接口
)

// Router 创建类型指定 Find 方法
type Router struct {
}

// ProtoParam 这里包含了 websocket 中的 sessionId 请求体和响应体
type ProtoParam struct {
	Request   *protobuf.Request
	Response  *protobuf.Response
	SessionID string
	DB        *gorm.DB
}

type fun func(*ProtoParam)

// Call 定义一个接口
type Call interface {
	call(*ProtoParam)
}

func (f fun) call(param *ProtoParam) {
	f(param)
}

// Find websocket 处理路由的主要方法
func (r *Router) Find(p *ProtoParam, f func(param *ProtoParam)) {
	handlerFind(p, fun(f))
}

func handlerFind(p *ProtoParam, c Call) {
	c.call(p)
}

// Handler 判断 websocket 传递的路由然后开始处理
func Handler(p *ProtoParam) {
	defaultRes := protobuf.Response{
		Code:   http.StatusInternalServerError,
		Msg:    "Internal Server Error",
		Type:   3,
		Data:   nil,
		Render: false,
		Html: &protobuf.Render{
			Code:   render.GetLayer(0, render.Incorrect, "Error", "Internal Server Error"),
			Type:   1,
			Id:     "layerMsgBox",
			Iframe: false,
		},
	}
	switch p.Request.Path {
	case APILogin:
		login := request.LoginParam{}
		err := json.Unmarshal(p.Request.Data, &login)
		if err != nil {
			defaultRes.Msg = err.Error()
			defaultRes.Html.Code = render.GetLayer(0, render.Incorrect, "Error", err.Error())
			p.Response = &defaultRes
		}
		data, errMsg, err := login.Exec(p.DB, p.SessionID)
		p.Response.Html.Code = render.GetLayer(0, render.Sad, "Login", errMsg)
		if err == nil {
			p.Response.Code = http.StatusOK
			p.Response.Html.Code = render.GetLayer(0, render.Smile, "Login", errMsg)
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
	}
}
