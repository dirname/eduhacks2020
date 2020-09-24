package control

import (
	"eduhacks2020/Go/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// DevicesControl 设备管理
type DevicesControl struct {
	Mongo *database.MongoClientDevice
}

type deviceResponse struct {
	Count int                     `json:"count"`
	Code  int                     `json:"code"`
	Msg   string                  `json:"msg"`
	Data  []database.ClientDevice `json:"data"`
}

// GetDevices 获取设备
func (d *DevicesControl) GetDevices(c *gin.Context) {
	var res []database.ClientDevice
	page, _ := strconv.Atoi(c.Request.FormValue("page"))
	limit, _ := strconv.Atoi(c.Request.FormValue("limit"))
	query := bson.M{}
	if nickname := c.Request.FormValue("nickname"); nickname != "" {
		query["name"] = bson.M{"$regex": nickname}
	}
	if username := c.Request.FormValue("username"); username != "" {
		query["user"] = bson.M{"$regex": username}
	}
	if systemID := c.Request.FormValue("systemId"); systemID != "" {
		query["systemId"] = bson.M{"$regex": systemID}
	}
	offset := (page - 1) * limit
	r := deviceResponse{}
	r.Code = 0
	msg := "OK"
	count, err := d.Mongo.Find(query, limit, offset, &res)
	if err != nil {
		msg = err.Error()
		r.Count = 0
		r.Code = http.StatusInternalServerError
		r.Data = nil
	}
	r.Count = count
	r.Data = res
	r.Msg = msg
	c.JSON(200, &r)
}
