package database

import (
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/utils"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

// ClientDevice 客户端的信息集合
type ClientDevice struct {
	ClientID    string `json:"clientId" bson:"clientId"`       // 标识ID
	SystemID    string `json:"systemId" bson:"systemId"`       // 系统ID
	ConnectTime uint64 `json:"connectTime" bson:"connectTime"` // 首次连接时间
	IsDeleted   bool   `json:"deleted" bson:"deleted"`         // 是否删除或下线
	UserName    string `json:"user" bson:"user"`               // 业务端标识用户账号
	NickName    string `json:"name" bson:"name"`               // 业务端标识用户昵称
	UserRole    int    `json:"role" bson:"role"`               // 业务端标识用户角色
}

// MongoClientDevice Mongo 来保存 ws 客户端的信息
type MongoClientDevice struct {
	Session        *mgo.Session
	CollectionName string
}

// Init 初始化
func (m *MongoClientDevice) Init() {
	session, err := mgo.DialWithInfo(setting.DialInfo)
	if err != nil {
		log.Errorf("failed to open mongo: %s", err.Error())
	}
	m.Session = session
	m.CollectionName = utils.GenUUIDv4()
	log.Infof("save the client device info to the collection %s", m.CollectionName)
}

// Close 删除集合
func (m *MongoClientDevice) Close() {
	if err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).DropCollection(); err != nil {
		log.WithFields(log.Fields{
			"reason": err.Error(),
		}).Error("drop collection failed")
	}
	defer m.Session.Close()
}

// Online 添加
func (m *MongoClientDevice) Online(d *ClientDevice) {
	if count, err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Find(bson.M{"systemId": d.SystemID}).Count(); err == nil {
		if count > 0 {
			findRes := ClientDevice{}
			m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Find(bson.M{"systemId": d.SystemID}).One(&findRes)
			d.UserRole = findRes.UserRole
			d.UserName = findRes.UserName
			d.NickName = findRes.NickName
			log.Infof("the same systemID has been connected and will be updated in the database systemID: %s", d.SystemID)
			selector := bson.M{"systemId": d.SystemID}
			update := bson.M{"$set": d}
			if err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Update(selector, update); err != nil {
				log.WithFields(log.Fields{
					"reason": err.Error(),
				}).Error("client failed to update")
			}
			return
		}
	}
	if err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Insert(d); err != nil {
		log.WithFields(log.Fields{
			"reason": err.Error(),
		}).Error("client failed to add to collection")
	}
}

// Offline 下线
func (m *MongoClientDevice) Offline(clientID string) {
	if err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Remove(bson.M{"clientId": clientID}); err != nil {
		log.WithFields(log.Fields{
			"reason": err.Error(),
		}).Error("client failed to remove collection")
	}
}

// SetUser 更新用户信息
func (m *MongoClientDevice) SetUser(clientID, user, name string, role int) {
	selector := bson.M{"clientId": clientID}
	update := bson.M{"$set": bson.M{"user": user, "name": name, "role": role}}
	if err := m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Update(selector, update); err != nil {
		log.WithFields(log.Fields{
			"reason": err.Error(),
		}).Error("client failed to set user")
	}
}

// Find 列出客户端
func (m *MongoClientDevice) Find(query interface{}, limit, offset int, c *[]ClientDevice) (int, error) {
	m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Find(query).Skip(offset).Limit(limit).All(c)
	return m.Session.DB(setting.DialInfo.Source).C(m.CollectionName).Find(query).Count()
}
