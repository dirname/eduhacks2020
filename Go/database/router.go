package database

import (
	"bytes"
	"crypto/md5"
	"eduhacks2020/Go/api/college"
	"eduhacks2020/Go/api/users"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	"encoding/hex"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
)

// 定义一些常量错误
const (
	signInvalid = "data has been tampered with invalid sign"
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
	Redis     *redis.Client
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

// 用于计算签名
func calcSign(timestamp string, data []byte) string {
	var buffer bytes.Buffer
	buffer.Write([]byte(timestamp))
	buffer.Write(data)
	h := md5.New()
	h.Write(buffer.Bytes())
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 返回校验的结果
func verifySign(salt string, submitSign, data []byte) bool {
	submit := string(submitSign)
	calc := calcSign(salt, data)
	return submit == calc
}

// Handler 判断 websocket 传递的路由然后开始处理
func Handler(p *ProtoParam) {
	switch p.Request.Path {
	case APILogin:
		login := users.LoginParam{}
		if err := json.Unmarshal(p.Request.Data, &login); err != nil {
			p.Response.Msg = err.Error()
			p.Response.Html.Code = render.GetLayer(0, render.Incorrect, "Error", err.Error())
			return
		}
		if !verifySign(login.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			p.Response.Html.Code = render.GetLayer(0, render.Sad, "Error", signInvalid)
			return
		}
		data, errMsg, err := login.Exec(p.DB, p.Redis, p.SessionID)
		p.Response.Html.Code = render.GetLayer(0, render.Sad, "Login", errMsg)
		if err == nil {
			p.Response.Code = http.StatusOK
			p.Response.Html.Code = render.GetLayer(0, render.Smile, "Login", errMsg)
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	case APIManagerStudentGet:
		p.Response.Html = nil
		p.Response.Render = false
		p.Response.Type = 5
		get := users.StudentGetParam{}
		if err := json.Unmarshal(p.Request.Data, &get); err != nil {
			p.Response.Msg = err.Error()
			return
		}
		if !verifySign(get.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			return
		}
		data, errMsg, err := get.Exec(p.DB, p.Redis)
		if err != nil {
			if err.Error() == users.TokenInvalid {
				p.Response.Code = -1
			}
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	case APIManagerCollegeGet:
		p.Response.Html = nil
		p.Response.Render = false
		p.Response.Type = 5
		get := college.GetParam{}
		if err := json.Unmarshal(p.Request.Data, &get); err != nil {
			p.Response.Msg = err.Error()
			return
		}
		if !verifySign(get.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			return
		}
		data, errMsg, err := get.Exec(p.DB, p.Redis)
		if err != nil {
			if err.Error() == users.TokenInvalid {
				p.Response.Code = -1
			}
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	case APIManagerCollegeAdd:
		add := college.AddParam{}
		if err := json.Unmarshal(p.Request.Data, &add); err != nil {
			p.Response.Msg = err.Error()
			p.Response.Html.Code = render.GetMsg(err.Error(), 3)
			return
		}
		if !verifySign(add.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			p.Response.Html.Code = render.GetMsg(signInvalid, 3)
			return
		}
		data, errMsg, err := add.Exec(p.DB, p.Redis)
		p.Response.Html.Code = render.GetMsg(errMsg, 3)
		if err == nil {
			p.Response.Code = http.StatusOK
			p.Response.Html.Code = render.GetMsg(errMsg, 3)
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	case APIManagerCollegeDelete:
		del := college.DelParam{}
		if err := json.Unmarshal(p.Request.Data, &del); err != nil {
			p.Response.Msg = err.Error()
			p.Response.Html.Code = render.GetMsg(err.Error(), 3)
			return
		}
		if !verifySign(del.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			p.Response.Html.Code = render.GetMsg(signInvalid, 3)
			return
		}
		data, errMsg, err := del.Exec(p.DB, p.Redis)
		p.Response.Html.Code = render.GetMsg(errMsg, 3)
		if err == nil {
			p.Response.Code = http.StatusOK
			p.Response.Html.Code = render.GetMsg(errMsg, 3)
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	case APILogout:
		logout := users.LogoutParam{}
		if err := json.Unmarshal(p.Request.Data, &logout); err != nil {
			p.Response.Msg = err.Error()
			p.Response.Html.Code = render.GetLayer(0, render.Incorrect, "Error", err.Error())
			return
		}
		if !verifySign(logout.Salt, p.Request.Sign, p.Request.Data) {
			p.Response.Msg = signInvalid
			p.Response.Html.Code = render.GetLayer(0, render.Sad, "Error", signInvalid)
			return
		}
		data, errMsg, err := logout.Exec(p.Redis)
		p.Response.Html.Code = render.GetLayer(0, render.Sad, "Logout", errMsg)
		if err == nil {
			p.Response.Code = http.StatusOK
			p.Response.Html.Code = render.GetLayer(0, render.Smile, "Logout", errMsg)
		}
		p.Response.Data = data
		p.Response.Msg = errMsg
		p.Response.Id = p.Request.Id
	}
}
