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

// MajorGetParam 专业获取的数据结构
type MajorGetParam struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Salt      string `json:"salt"`
	Token     string `json:"token"`
	Name      string `json:"majorName"`
	CollegeID uint   `json:"collegeID"`
}

// MajorGetView 获取专业的html渲染
type MajorGetView struct {
	Salt      string `json:"salt"`
	Token     string `json:"token"`
	CollegeID uint   `json:"collegeID"`
}

// MajorViewRes 专业的 html 响应
type MajorViewRes struct {
	ID        uint   `json:"id"`
	MajorID   string `json:"mid"`
	MajorName string `json:"name"`
}

// MajorResInfo 返回专业的数据结构
type MajorResInfo struct {
	ID          uint      `json:"id"`
	MajorID     string    `json:"mid"`
	MajorName   string    `json:"name"`
	CollegeID   uint      `json:"collegeID"`
	CollegeName string    `json:"collegeName"`
	CreatedAt   time.Time `json:"create"`
}

// MajorAddParam 添加专业的请求参数
type MajorAddParam struct {
	Salt  string `json:"salt"`
	Name  string `json:"name"`
	ID    uint   `json:"id"`
	Token string `json:"token"`
}

// MajorDelParam 删除专业的请求参数
type MajorDelParam struct {
	Salt  string `json:"salt"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// UpdateMajorParam 修改专业的请求参数
type UpdateMajorParam struct {
	Salt    string `json:"salt"`
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	College uint   `json:"collegeID"`
	Token   string `json:"token"`
}

// Exec 执行更新
func (c *UpdateMajorParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			db.Model(&psql.Major{}).Where("id = ?", c.ID).Updates(psql.Major{
				MajorName: c.Name,
				CollegeID: c.College,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行删除
func (d *MajorDelParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			db.Where("id = ?", d.ID).Delete(&psql.Major{})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行添加
func (c *MajorAddParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
		errMsg = users.TokenInvalid
		r.Code = -1
		r.Data = nullJs
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			errMsg = users.TokenInvalid
			r.Code = -1
			r.Data = nullJs
		} else {
			db.Create(&psql.Major{
				CollegeID: c.ID,
				MajorName: c.Name,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行查询
func (c *MajorGetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Html = nil
	r.Render = false
	r.Type = 5
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, &c); err != nil {
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
	errMsg := "OK"
	r.Data = nullJs
	r.Code = http.StatusOK
	if claims, err := utils.ParseToken(c.Token); err != nil {
		errMsg = users.TokenInvalid
		r.Code = -1
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			errMsg = users.TokenInvalid
			r.Code = -1
		} else {
			var collegeRows []MajorResInfo
			result := db.Model(&psql.Major{}).Select("majors.*, colleges.college_name, colleges.deleted_at").Joins("LEFT JOIN college.colleges on majors.college_id = colleges.id").Where(&psql.Major{
				CollegeID: c.CollegeID,
			}).Where("majors.major_name LIKE ?", fmt.Sprintf("%%%s%%", c.Name)).Where("colleges.deleted_at is null").Find(&collegeRows).Offset(c.Limit * (c.Page - 1)).Limit(c.Limit)
			res := response.TableResponse{
				Code:    0,
				Data:    collegeRows,
				Message: "OK",
				Count:   result.RowsAffected,
			}
			js, err := json.Marshal(&res)
			if err != nil {
				r.Data = nil
				r.Code = http.StatusInternalServerError
				errMsg = err.Error()
			}
			r.Data = js
		}
	}
	r.Msg = errMsg
}

// Exec 获取 html 的渲染
func (c *MajorGetView) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			var majorRows []MajorViewRes
			result := db.Model(&psql.Major{}).Where("college_id = ?", c.CollegeID).Find(&majorRows)
			html := ""
			for _, row := range majorRows {
				html += fmt.Sprintf("<option value=\"%d\">%s</option>\n", row.ID, row.MajorName)
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
