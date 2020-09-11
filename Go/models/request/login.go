package request

import (
	"eduhacks2020/Go/crypto"
	"eduhacks2020/Go/models/psql"
	"eduhacks2020/Go/models/response"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
	"time"
)

const (
	passwordValid = "password is mismatch"
	userBanned    = "account blocked"
)

// 登录使用的参数
type LoginParam struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}

// 登录后的结果
type LoginResponse struct {
	Token string               `json:"token"`
	Data  response.StudentInfo `json:"data"`
	Time  time.Time            `json:"time"`
}

// Exec 执行登录
func (l *LoginParam) Exec(db *gorm.DB, id string) ([]byte, string, error) {
	cipher := crypto.ChaCha20Poly1305{}
	cipher.Init()
	//var result interface{}
	switch l.Type {
	case 2:
		result := psql.Student{}
		row := db.Where("username = ?", l.Username).Or("phone = ?", l.Username).First(&result)
		if row.Error != nil {
			if row.Error == gorm.ErrRecordNotFound {
				return nil, "username or password is invalid", row.Error
			}
			return nil, "unknown error: " + row.Error.Error(), row.Error
		}
		restorePwd, _ := cipher.DecryptedFromHex(result.Password)
		if string(restorePwd) != l.Password {
			return nil, "username or password is invalid", errors.New(passwordValid)
		} else {
			if result.Banned {
				return nil, "this account is banned", errors.New(userBanned)
			}
			//generate token
			claims := utils.CustomClaims{
				UID:      result.UserID,
				Name:     result.Nickname,
				Username: result.Username,
				Phone:    result.Phone,
				Role:     2,
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix() - 1000, // 签名生效时间
					ExpiresAt: time.Now().Unix() + 3600, // 过期时间 一小时
					Issuer:    utils.Issuer,             //签名的发行者
				},
			}
			token, err := utils.GenerateToken(claims)
			if err != nil {
				return nil, err.Error(), err
			}
			studentInfo := response.StudentInfo{}
			db.Model(&psql.Student{}).Select("student.users.*,college.classes.class_name,college.classes.class_id,college.majors.major_name,college.majors.major_id,college.colleges.college_name,college.colleges.college_id").
				Joins("left join college.classes on student.users.class_id = college.classes.id left join college.majors on college.classes.major_id = college.majors.id LEFT JOIN college.colleges on college.majors.college_id = college.colleges.id").Scan(&studentInfo)
			res := LoginResponse{
				Token: token,
				Data:  studentInfo,
				Time:  time.Now(),
			}
			js, err := json.Marshal(&res)
			if err != nil {
				return nil, err.Error(), err
			}
			//session := database.SessionManager{Values: make(map[interface{}]interface{})}
			//session.Values["login"] = true
			//cipherText, err := session.EncryptedData(database.SessionName)
			//if err != nil {
			//	log.Errorf("Encrypted an error has occurred %s", err.Error())
			//}
			//session.SaveData(id, cipherText)
			return js, "Login success !", nil
		}
	}
	return nil, "", nil
}
