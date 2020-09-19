package users

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/models/psql"
	"eduhacks2020/Go/models/response"
	"eduhacks2020/Go/protobuf"
	"eduhacks2020/Go/utils"
	"encoding/json"
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
	}
	redisAuth := api.AuthRedis{Redis: redis}
	flag, _ := redisAuth.GetFlag(claims.UID)
	if claims.Flag != flag {
		r.Data = nullJs
		r.Code = -1
		errMsg = TokenInvalid
	} else {
		var stuRows []response.StudentInfo
		resTemp := psql.Student{}
		result := db.Where(&psql.Student{
			Username: l.Username,
			Phone:    l.Phone,
			Email:    l.Email,
		}).Find(&resTemp)
		db.Model(&psql.Student{}).Select("student.users.*,college.classes.deleted_at,college.classes.class_name,college.classes.class_id,college.majors.deleted_at,college.majors.major_name,college.majors.major_id,college.colleges.deleted_at,college.colleges.college_name,college.colleges.college_id").
			Joins("left join college.classes on student.users.class_id = college.classes.id left join college.majors on college.classes.major_id = college.majors.id LEFT JOIN college.colleges on college.majors.college_id = college.colleges.id").
			Where(&psql.Student{
				Username: l.Username,
				Phone:    l.Phone,
				Email:    l.Email,
				Nickname: l.Nickname,
			}).
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
	r.Msg = errMsg
}
