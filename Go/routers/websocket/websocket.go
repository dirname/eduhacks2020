package websocket

import (
	"eduhacks2020/Go/api/college"
	"eduhacks2020/Go/api/users"
	"eduhacks2020/Go/protobuf"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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

// Handler 判断 websocket 传递的路由然后开始处理
func Handler(p *ProtoParam) {
	switch p.Request.Path {
	case APILogin:
		login := users.LoginParam{}
		login.Exec(p.DB, p.Redis, p.SessionID, p.Request, p.Response)
	case APIManagerStudentGet:
		get := users.StudentGetParam{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerCollegeGet:
		get := college.GetParam{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerCollegeView:
		get := college.GetView{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerCollegeAdd:
		add := college.AddParam{}
		add.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerCollegeDelete:
		del := college.DelParam{}
		del.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerCollegeEdit:
		edit := college.UpdateParam{}
		edit.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerMajorGet:
		get := college.MajorGetParam{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerMajorView:
		get := college.MajorGetView{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerMajorAdd:
		add := college.MajorAddParam{}
		add.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerMajorDelete:
		del := college.MajorDelParam{}
		del.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerMajorEdit:
		edit := college.UpdateMajorParam{}
		edit.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerClassGet:
		get := college.ClassGetParam{}
		get.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerClassView:
		view := college.ClassGetView{}
		view.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerClassAdd:
		add := college.ClassAddParam{}
		add.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerClassDelete:
		del := college.ClassDeleteParam{}
		del.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APIManagerClassEdit:
		edit := college.ClassUpdateParam{}
		edit.Exec(p.DB, p.Redis, p.Request, p.Response)
	case APILogout:
		logout := users.LogoutParam{}
		logout.Exec(p.Redis, p.Request, p.Response)
	}
}
