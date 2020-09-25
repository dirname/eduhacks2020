package users

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/crypto"
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
)

// TokenInvalid 登录状态失效的错误信息
const TokenInvalid = "token is invalid, please login again"

// StudentGetParam 学生获取的请求结构
type StudentGetParam struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Salt     string `json:"salt"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

// StudentAddParam 学生添加的请求结构
type StudentAddParam struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	ClassID  uint   `json:"classID"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Salt     string `json:"salt"`
}

// StudentDelBanParam 学生封禁删除的请求参数
type StudentDelBanParam struct {
	Username string `json:"username"`
	Type     int    `json:"type"`
	Token    string `json:"token"`
	Salt     string `json:"salt"`
}

// TeacherGetParam 教师获取的
type TeacherGetParam struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Salt     string `json:"salt"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}

// AcademicGetParam 教务获取的
type AcademicGetParam struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Salt     string `json:"salt"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Token    string `json:"token"`
}

// TeacherAddParam 教师添加的
type TeacherAddParam struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Gender   string `json:"gender"`
	FullTime string `json:"fullTime"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Salt     string `json:"salt"`
}

// TeacherDelParam 教师删除的
type TeacherDelParam struct {
	Username string `json:"username"`
	Type     int    `json:"type"`
	Token    string `json:"token"`
	Salt     string `json:"salt"`
}

// Exec 执行查询
func (l *StudentGetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Html = nil
	r.Render = false
	r.Type = 5
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		r.Msg = err.Error()
		return
	}
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
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
	claims, err := utils.ParseToken(l.Token)
	if err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = TokenInvalid
		} else {
			var stuRows []response.StudentInfo
			resTemp := psql.Student{}
			result := db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", l.Username)).
				Where("phone LIKE ?", fmt.Sprintf("%%%s%%", l.Phone)).
				Where("email LIKE ?", fmt.Sprintf("%%%s%%", l.Email)).
				Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", l.Nickname)).Find(&resTemp)
			db.Model(&psql.Student{}).Select("student.users.*,college.classes.deleted_at,college.classes.class_name,college.classes.class_id,college.majors.deleted_at,college.majors.major_name,college.majors.major_id,college.colleges.deleted_at,college.colleges.college_name,college.colleges.college_id").
				Joins("left join college.classes on student.users.class_id = college.classes.id left join college.majors on college.classes.major_id = college.majors.id LEFT JOIN college.colleges on college.majors.college_id = college.colleges.id").
				//Where(&psql.Student{
				//	Username: l.Username,
				//	Phone:    l.Phone,
				//	Email:    l.Email,
				//	Nickname: l.Nickname,
				//}).
				Where("student.users.username LIKE ?", fmt.Sprintf("%%%s%%", l.Username)).
				Where("student.users.phone LIKE ?", fmt.Sprintf("%%%s%%", l.Phone)).
				Where("student.users.email LIKE ?", fmt.Sprintf("%%%s%%", l.Email)).
				Where("student.users.nickname LIKE ?", fmt.Sprintf("%%%s%%", l.Nickname)).
				Where("college.colleges.deleted_at is null").
				Where("college.majors.deleted_at is null").
				Where("college.classes.deleted_at is null").Offset((l.Page - 1) * l.Limit).Limit(l.Limit).Find(&stuRows)
			res := response.TableResponse{
				Code:    0,
				Data:    stuRows,
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

// Exec 执行查询
func (t *TeacherGetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Html = nil
	r.Render = false
	r.Type = 5
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, t); err != nil {
		r.Msg = err.Error()
		return
	}
	if !utils.VerifySign(t.Salt, request.Sign, request.Data) {
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
	claims, err := utils.ParseToken(t.Token)
	if err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = TokenInvalid
		} else {
			var stuRows []response.TeacherInfo
			resTemp := psql.Teacher{}
			result := db.Model(&psql.Teacher{}).Where("username LIKE ?", fmt.Sprintf("%%%s%%", t.Username)).
				Where("phone LIKE ?", fmt.Sprintf("%%%s%%", t.Phone)).
				Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", t.Nickname)).Find(&resTemp)
			db.Model(&psql.Teacher{}).Where("username LIKE ?", fmt.Sprintf("%%%s%%", t.Username)).
				Where("phone LIKE ?", fmt.Sprintf("%%%s%%", t.Phone)).
				Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", t.Nickname)).
				Offset((t.Page - 1) * t.Limit).Limit(t.Limit).Find(&stuRows)
			res := response.TableResponse{
				Code:    0,
				Data:    stuRows,
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

// Exec 执行添加
func (l *StudentAddParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
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
	if claims, err := utils.ParseToken(l.Token); err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = TokenInvalid
		} else {
			cipher := crypto.ChaCha20Poly1305{}
			cipher.Init()
			gender := false
			switch l.Gender {
			case "false":
				gender = false
			case "true":
				gender = true
			default:
				gender = false
			}
			password := cipher.EncryptedToHex(l.Password)
			db.Create(&psql.Student{
				Username: l.Username,
				Nickname: l.Nickname,
				Gender:   gender,
				Phone:    l.Phone,
				Email:    l.Email,
				ClassID:  l.ClassID,
				Password: password,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行删除或添加
func (l *StudentDelBanParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
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
	if claims, err := utils.ParseToken(l.Token); err != nil {
		r.Data = nullJs
		errMsg = TokenInvalid
		r.Code = -1
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			errMsg = TokenInvalid
			r.Code = -1
		} else {
			switch l.Type {
			case 0:
				db.Where("username = ?", l.Username).Delete(&psql.Student{})
				errMsg = "Deleted successfully !"
			case 1:
				db.Model(&psql.Student{}).Where("username = ?", l.Username).Updates(psql.Student{
					Banned: true,
				})
				errMsg = "Banned successfully !"
			case 2:
				db.Model(&psql.Student{}).Where("username = ?", l.Username).Update("banned", false)
				errMsg = "Unblocked successfully !"
			}
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行添加
func (l *TeacherAddParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
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
	if claims, err := utils.ParseToken(l.Token); err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = TokenInvalid
		} else {
			cipher := crypto.ChaCha20Poly1305{}
			cipher.Init()
			gender := false
			switch l.Gender {
			case "false":
				gender = false
			case "true":
				gender = true
			default:
				gender = false
			}
			fullTime := false
			switch l.FullTime {
			case "false":
				fullTime = false
			case "true":
				fullTime = true
			default:
				fullTime = false
			}
			password := cipher.EncryptedToHex(l.Password)
			db.Create(&psql.Teacher{
				Username: l.Username,
				Nickname: l.Nickname,
				Gender:   gender,
				Phone:    l.Phone,
				Email:    l.Email,
				FullTime: fullTime,
				Password: password,
			})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行删除
func (l *TeacherDelParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, l); err != nil {
		r.Msg = err.Error()
		r.Html.Code = render.GetMsg(err.Error(), 3)
		return
	}
	if !utils.VerifySign(l.Salt, request.Sign, request.Data) {
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
	errMsg := "Deleted successfully !"
	if claims, err := utils.ParseToken(l.Token); err != nil {
		r.Data = nullJs
		errMsg = TokenInvalid
		r.Code = -1
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			errMsg = TokenInvalid
			r.Code = -1
		} else {
			db.Where("username = ?", l.Username).Delete(&psql.Teacher{})
		}
	}
	r.Html.Code = render.GetMsg(errMsg, 3)
	r.Msg = errMsg
}

// Exec 执行查询
func (t *AcademicGetParam) Exec(db *gorm.DB, redis *redis.Client, request *protobuf.Request, r *protobuf.Response) {
	r.Html = nil
	r.Render = false
	r.Type = 5
	r.Id = request.Id
	if err := json.Unmarshal(request.Data, t); err != nil {
		r.Msg = err.Error()
		return
	}
	if !utils.VerifySign(t.Salt, request.Sign, request.Data) {
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
	claims, err := utils.ParseToken(t.Token)
	if err != nil {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			r.Data = nullJs
			r.Code = -1
			errMsg = TokenInvalid
		} else {
			var stuRows []response.TeacherInfo
			resTemp := psql.Teacher{}
			result := db.Model(&psql.Teacher{}).Where("username LIKE ?", fmt.Sprintf("%%%s%%", t.Username)).
				Where("phone LIKE ?", fmt.Sprintf("%%%s%%", t.Phone)).
				Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", t.Nickname)).Find(&resTemp)
			db.Model(&psql.Teacher{}).Where("username LIKE ?", fmt.Sprintf("%%%s%%", t.Username)).
				Where("phone LIKE ?", fmt.Sprintf("%%%s%%", t.Phone)).
				Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", t.Nickname)).
				Offset((t.Page - 1) * t.Limit).Limit(t.Limit).Find(&stuRows)
			res := response.TableResponse{
				Code:    0,
				Data:    stuRows,
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
