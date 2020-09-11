package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些赋值
var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("that's not even a token")
	ErrTokenInvalid     = errors.New("couldn't handle this token ")
	SignKey             = "secretJWTToken"
	Issuer              = "JWTTokenManager"
)

// CustomClaims 载荷, token 里的信息
type CustomClaims struct {
	UID      string `json:"userId"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Flag     string `json:"flag"`
	Role     int    `json:"role"`
	jwt.StandardClaims
}

// NewJWT() 新建一个 JWT 实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(SignKey),
	}
}

// 获取 signKey
func (j *JWT) getSignKey() []byte {
	return j.SigningKey
}

// 设置 signKey
func (j *JWT) setSignKey(key string) {
	j.SigningKey = []byte(key)
}

// CreateToken 生成一个token
func (j *JWT) createToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 Token
func (j *JWT) parseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

// 更新token
func (j *JWT) refreshToken(tokenString string) (string, error) {
	//jwt.TimeFunc = func() time.Time {
	//	return time.Unix(0, 0)
	//}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return "", ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return "", ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return "", ErrTokenNotValidYet
			} else {
				return "", ErrTokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.createToken(*claims)
	}
	return "", ErrTokenInvalid
}

// GenerateToken 生成一个令牌
func GenerateToken(claims CustomClaims) (string, error) {
	j := NewJWT()
	token, err := j.createToken(claims)
	return token, err
}
