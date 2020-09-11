package database

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/securecookie"
	"github.com/kidstuff/mongostore"
	log "github.com/sirupsen/logrus"
	"time"
)

const sessionKey = "EduHacks2020.*" //session 默认使用的密匙
const sessionDB = "sessions"        //session 默认使用的数据库
const SessionName = "token"         //session 默认的名称

var codecs = securecookie.CodecsFromPairs([]byte(sessionKey))

var dialInfo = &mgo.DialInfo{
	Addrs:    []string{mongoHost},
	Source:   mongoDB,
	Username: mongoUser,
	Password: mongoPwd,
}

// SessionManager 这是对 mongoStore 的改写, 从 mgo 层面直接解密 session, 使其在 websocket 通信时能实时读取
type SessionManager struct {
	Values map[interface{}]interface{}
}

// CreateMongoStore 返回 *mgo.Session 再调用结束后及时释放数据库连接资源
func CreateMongoStore() (*mongostore.MongoStore, *mgo.Session) {
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Error(err)
	}
	return mongostore.NewMongoStore(session.DB(mongoDB).C(sessionDB), 86400, true,
		[]byte(sessionKey)), session
}

// GetData 对 mongoStore 的附加部分, 从 mgo 直接读取 session.id 中的数据, 随后加解密
func (*SessionManager) GetData(id string) (interface{}, error) {
	session, err := mgo.DialWithInfo(dialInfo)
	defer session.Close()
	if err != nil {
		log.Error(err)
	}
	objectId := bson.ObjectIdHex(id)
	c := session.DB(mongoDB).C(sessionDB)
	var one map[string]interface{}
	err = c.FindId(objectId).One(&one)
	return one["data"], err
}

// DecryptedData 解密 session 的数据
func (s *SessionManager) DecryptedData(data string, sessionName string) {
	defer securecookie.DecodeMulti(sessionName, data, &s.Values,
		codecs...)
}

// EncryptedData 加密 session 的数据
func (s *SessionManager) EncryptedData(sessionName string) (string, error) {
	encoded, err := securecookie.EncodeMulti(sessionName, s.Values,
		codecs...)
	return encoded, err
}

// SaveData 保存 session 的数据
func (s *SessionManager) SaveData(id string, data string) error {
	session, err := mgo.DialWithInfo(dialInfo)
	defer session.Close()
	if err != nil {
		log.Error(err)
	}
	objectId := bson.ObjectIdHex(id)
	c := session.DB(mongoDB).C(sessionDB)
	err = c.UpdateId(objectId, bson.M{"data": data, "modified": time.Now()})
	return err
}

// DeleteData 删除 session 的数据
func (s *SessionManager) DeleteData(id string) error {
	session, err := mgo.DialWithInfo(dialInfo)
	defer session.Close()
	if err != nil {
		log.Error(err)
	}
	objectId := bson.ObjectIdHex(id)
	c := session.DB(mongoDB).C(sessionDB)
	err = c.RemoveId(objectId)
	return err
}
