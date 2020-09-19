package college

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/api/users"
	"eduhacks2020/Go/models/psql"
	"eduhacks2020/Go/models/response"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/render"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// GetParam 获取学院的请求参数
type GetParam struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Salt  string `json:"salt"`
	Token string `json:"token"`
}

// GetView 获取学院的html渲染
type GetView struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Salt  string `json:"salt"`
	Token string `json:"token"`
}

// ResInfo 返回学院的数据结构
type ResInfo struct {
	ID          uint      `json:"id"`
	CollegeID   string    `json:"cid"`
	CollegeName string    `json:"name"`
	CreatedAt   time.Time `json:"create"`
}

// AddParam 添加学院的请求参数
type AddParam struct {
	Salt  string `json:"salt"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// UpdateParam 修改学院的请求参数
type UpdateParam struct {
	Salt  string `json:"salt"`
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// DelParam 删除学院的请求参数
type DelParam struct {
	Salt  string `json:"salt"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// Exec 执行删除
func (d *DelParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, d); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(d.Salt, request.Sign, request.Data) {
		r.Msg = utils.SignInvalid
		r.Html.Code = render.GetMsg(utils.SignInvalid, 3)
		return
	}
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	errMsg := "Deleted successfully !"
	r.Code = http.StatusOK
	r.Data = nil
	if claims, err := utils.ParseToken(d.Token); err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = users.TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = users.TokenInvalid
		} else {
			db.Where("id = ?", d.ID).Delete(&psql.College{})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行添加
func (c *AddParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, c); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(c.Salt, request.Sign, request.Data) {
		r.Msg = utils.SignInvalid
		r.Html.Code = render.GetMsg(utils.SignInvalid, 3)
		return
	}
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	r.Code = http.StatusOK
	r.Data = nil
	errMsg := "Added successfully !"
	if claims, err := utils.ParseToken(c.Token); err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = users.TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = users.TokenInvalid
		} else {
			db.Create(&psql.College{
				CollegeName: c.Name,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行更新
func (c *UpdateParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, c); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(c.Salt, request.Sign, request.Data) {
		r.Msg = utils.SignInvalid
		r.Html.Code = render.GetMsg(utils.SignInvalid, 3)
		return
	}
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	r.Code = http.StatusOK
	r.Data = nil
	errMsg := "Updated successfully !"
	if claims, err := utils.ParseToken(c.Token); err != nil {
		r.Data = nullJs
		errMsg = users.TokenInvalid
		r.Code = -1
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			errMsg = users.TokenInvalid
			r.Code = -1
		} else {
			db.Model(&psql.College{}).Where("id = ?", c.ID).Update("college_name", c.Name)
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行查询
func (c *GetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	r.Html = nil
	r.Render = false
	r.Type = 5
	if err := json.Unmarshal(request.Data, c); err != nil {
		r.Msg = err.Error()
		return
	}
	if !utils.VerifySign(c.Salt, request.Sign, request.Data) {
		r.Msg = utils.SignInvalid
		return
	}
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	r.Code = http.StatusOK
	r.Data = nil
	errMsg := "OK"
	if claims, err := utils.ParseToken(c.Token); err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = users.TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = users.TokenInvalid
		} else {
			var collegeRows []ResInfo
			result := db.Model(&psql.College{}).Find(&collegeRows)
			res := response.TableResponse{
				Code:    0,
				Data:    collegeRows,
				Message: "OK",
				Count:   result.RowsAffected,
			}
			js, err := json.Marshal(&res)
			if err != nil {
				errMsg = err.Error()
			}
			r.Data = js
		}
	}
	r.Msg = errMsg
}

// Exec 执行 html 的选择项
func (c *GetView) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Html = nil
	r.Render = false
	r.Type = 5
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, c); err != nil {
		r.Msg = err.Error()
		return
	}
	if !utils.VerifySign(c.Salt, request.Sign, request.Data) {
		r.Msg = utils.SignInvalid
		return
	}
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	r.Code = http.StatusOK
	r.Data = nil
	errMsg := "OK"
	if claims, err := utils.ParseToken(c.Token); err != nil {
		r.Code = -1
		r.Data = nullJs
		errMsg = users.TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Code = -1
			r.Data = nullJs
			errMsg = users.TokenInvalid
		} else {
			var collegeRows []ResInfo
			result := db.Model(&psql.College{}).Find(&collegeRows)
			html := ""
			for _, row := range collegeRows {
				html += fmt.Sprintf("<option value=\"%d\">%s</option>\n", row.ID, row.CollegeName)
			}
			res := response.TableResponse{
				Code:    0,
				Data:    html,
				Message: "OK",
				Count:   result.RowsAffected,
			}
			js, err := json.Marshal(&res)
			if err != nil {
				errMsg = err.Error()
			}
			r.Data = js
		}
	}
	r.Msg = errMsg
}
