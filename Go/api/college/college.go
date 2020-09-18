package college

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/api/users"
	"eduhacks2020/Go/models/psql"
	"eduhacks2020/Go/models/response"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

// GetParam 获取学院的请求参数
type GetParam struct {
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

// DelParam 删除学院的请求参数
type DelParam struct {
	Salt  string `json:"salt"`
	ID    int    `json:"name"`
	Token string `json:"token"`
}

// Exec 执行删除
func (d *DelParam) Exec(db *gorm.DB, redis *redis.Client) ([]byte, string, error) {
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	if claims, err := utils.ParseToken(d.Token); err != nil {
		return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
		}
	}
	db.Where("id = ?", d.ID).Delete(&psql.College{})
	return nil, "Deleted successfully !", nil
}

// Exec 执行添加
func (c *AddParam) Exec(db *gorm.DB, redis *redis.Client) ([]byte, string, error) {
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	if claims, err := utils.ParseToken(c.Token); err != nil {
		return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
		}
	}
	db.Create(&psql.College{
		CollegeName: c.Name,
	})
	return nil, "Added successfully !", nil
}

// Exec 执行查询
func (c *GetParam) Exec(db *gorm.DB, redis *redis.Client) ([]byte, string, error) {
	nullJs, _ := json.Marshal(response.TableResponse{
		Code:    -1,
		Data:    nil,
		Message: "",
		Count:   0,
	})
	if claims, err := utils.ParseToken(c.Token); err != nil {
		return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
	} else {
		redisAuth := api.AuthRedis{Redis: redis}
		flag, _ := redisAuth.GetFlag(claims.UID)
		if claims.Flag != flag {
			return nullJs, users.TokenInvalid, errors.New(users.TokenInvalid)
		}
	}

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

		return nil, err.Error(), err
	}
	return js, "ok", nil

}
