package users

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/crypto"
	"eduhacks2020/Go/models/psql"
	"eduhacks2020/Go/models/response"
	"eduhacks2020/Go/utils"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	passwordValid = "password is mismatch"
	userBanned    = "account blocked"
	unknownLogin  = "unknown login field"
)

// LoginParam 登录使用的参数
type LoginParam struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Type     int    `json:"type"`
	Salt     string `json:"salt"`
}

// LoginResponse 登录后的结果
type LoginResponse struct {
	Token string               `json:"token"`
	Data  response.StudentInfo `json:"data"`
	Time  time.Time            `json:"time"`
}

// Exec 执行登录
func (l *LoginParam) Exec(db *gorm.DB, redis *redis.Client, id string) ([]byte, string, error) {

	//var result interface{}
	switch l.Type {
	case 1:
		return l.teacherLogin(db, redis, id)
	case 2:
		return l.studentLogin(db, redis, id)
	default:
		return nil, "未知的登录域", errors.New(unknownLogin)
	}

}

// managerLogin 教务的登录

func (l *LoginParam) managerLogin(db *gorm.DB, redis *redis.Client, id string) ([]byte, string, error) {
	return nil, "", nil
}

// teacherLogin 教师的登录
func (l *LoginParam) teacherLogin(db *gorm.DB, redis *redis.Client, id string) ([]byte, string, error) {
	cipher := crypto.ChaCha20Poly1305{}
	cipher.Init()
	result := psql.Teacher{}
	row := db.Where("username = ?", l.Username).Or("phone = ?", l.Username).Or("email = ?", l.Username).First(&result)
	if row.Error != nil {
		if row.Error == gorm.ErrRecordNotFound {
			return nil, "username or password is invalid", row.Error
		}
		return nil, "unknown error: " + row.Error.Error(), row.Error
	}
	restorePwd, _ := cipher.DecryptedFromHex(result.Password)
	if string(restorePwd) != l.Password {
		return nil, "username or password is invalid", errors.New(passwordValid)
	}

	//generate token
	userFlag := utils.GenUUIDv5(result.UserID.String())
	claims := utils.CustomClaims{
		UID:      result.UserID.String(),
		Name:     result.Nickname,
		Username: result.Username,
		Phone:    result.Phone,
		Role:     1,
		Flag:     userFlag,
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
	return nil, token, nil
}

// studentLogin 学生的登录
func (l *LoginParam) studentLogin(db *gorm.DB, redis *redis.Client, id string) ([]byte, string, error) {
	cipher := crypto.ChaCha20Poly1305{}
	cipher.Init()
	result := psql.Student{}
	row := db.Where("username = ?", l.Username).Or("phone = ?", l.Username).Or("email = ?", l.Username).First(&result)
	if row.Error != nil {
		if row.Error == gorm.ErrRecordNotFound {
			return nil, "username or password is invalid", row.Error
		}
		return nil, "unknown error: " + row.Error.Error(), row.Error
	}
	restorePwd, _ := cipher.DecryptedFromHex(result.Password)
	if string(restorePwd) != l.Password {
		return nil, "username or password is invalid", errors.New(passwordValid)
	}
	if result.Banned {
		return nil, "this account is banned", errors.New(userBanned)
	}
	//generate token
	userFlag := utils.GenUUIDv5(result.UserID.String())
	claims := utils.CustomClaims{
		UID:      result.UserID.String(),
		Name:     result.Nickname,
		Username: result.Username,
		Phone:    result.Phone,
		Role:     2,
		Flag:     userFlag,
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
	redisAuth := api.AuthRedis{Redis: redis}
	err = redisAuth.SetFlag(claims.UID, userFlag)
	if err != nil {
		log.Error(err.Error())
	}
	return js, "Login success !", nil
}
