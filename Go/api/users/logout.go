package users

import (
	"context"
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/pkg/setting"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// LogoutParam 退出登录的请求参数
type LogoutParam struct {
	SystemID string `json:"systemId"`
	Token    string `json:"token"`
	Salt     string `json:"salt"`
}

// Exec 执行删除
func (l *LogoutParam) Exec(redis *redis.Client, request *protobuf.Request, response *protobuf.Response, id string, mongo *mgo.Session, collectionName string) {
	response.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		response.Msg = err.Error()
		response.Html.Code = render.GetLayer(0, render.Incorrect, "Error", err.Error())
		return
	}
	response.Id = request.Id
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
		response.Msg = utils.SignInvalid
		response.Html.Code = render.GetLayer(0, render.Sad, "Error", utils.SignInvalid)
		return
	}
	claims, err := utils.ParseToken(l.Token)
	errMsg := "Logout success !"
	if err != nil {
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			errMsg = TokenInvalid
		} else {
			redis.Del(context.Background(), claims.UID).Result()
		}
	}
	response.Html.Code = render.GetLayer(0, render.Sad, "Logout", errMsg)
	if err == nil {
		response.Code = http.StatusOK
		response.Html.Code = render.GetLayer(0, render.Smile, "Logout", errMsg)
		session := database.SessionManager{}
		session.DeleteData(id)
		selector := bson.M{"systemId": l.SystemID}
		update := bson.M{"$set": bson.M{"user": "", "name": "", "role": 0}}
		if err := mongo.DB(setting.DialInfo.Source).C(collectionName).Update(selector, update); err != nil {
			log.WithFields(log.Fields{
				"reason":   err.Error(),
				"systemID": l.SystemID,
			}).Error("client failed to log out user")
		}
	}
	response.Data = nil
	response.Msg = errMsg
}
