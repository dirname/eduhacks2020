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

// ClassGetParam 获取班级信息
type ClassGetParam struct {
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
	Salt    string `json:"salt"`
	Token   string `json:"token"`
	Name    string `json:"className"`
	MajorID uint   `json:"majorID"`
}

// ClassResInfo 班级的返回信息
type ClassResInfo struct {
	ID          uint      `json:"id"`
	ClassName   string    `json:"name"`
	MajorID     uint      `json:"majorID"`
	MajorName   string    `json:"majorName"`
	CollegeName string    `json:"collegeName"`
	CollegeID   uint      `json:"collegeID"`
	CreatedAt   time.Time `json:"create"`
}

// ClassGetView 获取班级的html渲染
type ClassGetView struct {
	Salt    string `json:"salt"`
	Token   string `json:"token"`
	MajorID uint   `json:"majorID"`
}

// ClassViewRes 班级的 html 响应
type ClassViewRes struct {
	ID        uint   `json:"id"`
	ClassID   string `json:"cid"`
	ClassName string `json:"name"`
}

// ClassAddParam 添加学院的信息
type ClassAddParam struct {
	Salt  string `json:"salt"`
	Name  string `json:"name"`
	Token string `json:"token"`
	Major uint   `json:"majorID"`
}

// ClassUpdateParam 修改学院的信息
type ClassUpdateParam struct {
	Salt    string `json:"salt"`
	Token   string `json:"token"`
	ID      uint   `json:"id"`
	Name    string `json:"className"`
	MajorID uint   `json:"majorID"`
}

// ClassDeleteParam 删除学院的信息
type ClassDeleteParam struct {
	Salt  string `json:"salt"`
	ID    int    `json:"id"`
	Token string `json:"token"`
}

// Exec 执行更新
func (c *ClassUpdateParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			db.Model(&psql.Class{}).Where("id = ?", c.ID).Updates(psql.Class{
				ClassName: c.Name,
				MajorID:   c.MajorID,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行删除
func (d *ClassDeleteParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			db.Where("id = ?", d.ID).Delete(&psql.Class{})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行查询
func (c *ClassGetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			var classRows []ClassResInfo
			result := db.Model(&psql.Class{}).Select("classes.*, majors.college_id, majors.major_name, majors.deleted_at, colleges.college_name, colleges.deleted_at").Joins("LEFT JOIN college.majors on classes.major_id = majors.id LEFT JOIN college.colleges on majors.college_id = colleges.id").Where(&psql.Class{
				MajorID: c.MajorID,
			}).Where("classes.class_name LIKE ?", fmt.Sprintf("%%%s%%", c.Name)).Where("colleges.deleted_at is null").Where("majors.deleted_at is null").Find(&classRows)
			res := response.TableResponse{
				Code:    0,
				Data:    classRows,
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

// Exec 执行添加
func (c *ClassAddParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			db.Create(&psql.Class{
				ClassName: c.Name,
				MajorID:   c.Major,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 获取 html 的渲染
func (c *ClassGetView) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
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
			var classRows []ClassViewRes
			result := db.Model(&psql.Major{}).Where("major_id = ?", c.MajorID).Find(&classRows)
			html := ""
			for _, row := range classRows {
				html += fmt.Sprintf("<option value=\"%d\">%s</option>\n", row.ID, row.ClassName)
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
