package database

import (
	"eduhacks2020/Go/pkg/setting"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"testing"
)

type devices struct {
	ClientID    string `json:"clientId" bson:"clientId"`       // 标识ID
	SystemID    string `json:"systemId" bson:"systemId"`       // 系统ID
	ConnectTime uint64 `json:"connectTime" bson:"connectTime"` // 首次连接时间
	IsDeleted   bool   `json:"deleted" bson:"deleted"`         // 是否删除或下线
	UserName    string `json:"user" bson:"user"`               // 业务端标识用户账号
	NickName    string `json:"name" bson:"name"`               // 业务端标识用户昵称
	UserRole    int    `json:"role" bson:"role"`               // 业务端标识用户角色
}

func TestRun(t *testing.T) {
	setting.DefaultSetting()
	session, err := mgo.DialWithInfo(setting.DialInfo)
	if err != nil {
		log.Errorf("Failed to open mongo: %s", err.Error())
	}
	device := devices{
		ClientID:    "123",
		SystemID:    "123",
		ConnectTime: 0,
		IsDeleted:   false,
		UserName:    "123",
		NickName:    "123",
		UserRole:    0,
	}
	//session.DB(setting.DialInfo.Source).C("devices").DropCollection()
	session.DB(setting.DialInfo.Source).C("devices").Insert(&device)
	defer session.Close()
}

func TestUpdate(t *testing.T) {
	setting.DefaultSetting()
	session, err := mgo.DialWithInfo(setting.DialInfo)
	if err != nil {
		log.Errorf("Failed to open mongo: %s", err.Error())
	}
	defer session.Close()
	session.DB(setting.DialInfo.Source).C("devices").Update(bson.M{"clientId": "123"}, bson.M{"$set": bson.M{"user": "456", "nickname": "456", "role": 8}})
}

func TestFind(t *testing.T) {
	setting.DefaultSetting()
	session, err := mgo.DialWithInfo(setting.DialInfo)
	if err != nil {
		log.Errorf("Failed to open mongo: %s", err.Error())
	}
	defer session.Close()
	var res []devices
	session.DB(setting.DialInfo.Source).C("devices").Find(bson.M{"clientId": bson.M{"$regex": "1"}}).All(&res)
	fmt.Printf("%v", res)
}
