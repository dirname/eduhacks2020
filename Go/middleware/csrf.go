package middleware

import (
	"eduhacks2020/Go/crypto"
	"eduhacks2020/Go/database"
	"eduhacks2020/Go/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type csrfToken struct {
	SessionId string    `json:"session_id"`
	MaxAge    time.Time `json:"max_age"`
}

// GenerateCSRF 生成 csrf token
func GenerateCSRF(id string) string {
	csrfToken := csrfToken{
		SessionId: id,
		MaxAge:    time.Now().Add(24 * time.Hour),
	}
	cipher := crypto.ChaCha20Poly1305{}
	cipher.Init()
	tokenObject, _ := json.Marshal(csrfToken)
	return cipher.EncryptedToBase64(string(tokenObject))
}

// CSRF gin 的中间件
func CSRF() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		clientIP := context.ClientIP()
		method := context.Request.Method

		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			context.Next()
			return
		}

		res := models.Response{
			Code:   http.StatusBadRequest,
			Path:   path,
			Method: method,
			Msg:    "Bad Request Invalid CSRF",
			Data:   nil,
			Time:   time.Now(),
			IP:     clientIP,
		}

		csrf := context.Request.FormValue("csrf")
		if csrf == "" {
			context.Abort()
			context.JSON(res.Code, res)
			return
		}
		store, dbSession := database.CreateMongoStore()
		defer dbSession.Close()
		session, err := store.Get(context.Request, database.SessionName)

		csrfToken := csrfToken{
			SessionId: session.ID,
			MaxAge:    time.Now().Add(24 * time.Hour),
		}

		if err != nil {
			context.Abort()
			context.JSON(res.Code, res)
			return
		}

		cipher := crypto.ChaCha20Poly1305{}
		cipher.Init()

		decode, err := cipher.DecryptedFromBase64(csrf)
		if err != nil {
			res.Msg = "Bad Request Got Invalid CSRF " + err.Error()
			context.Abort()
			context.JSON(res.Code, res)
			return
		}
		err = json.Unmarshal(decode, &csrfToken)
		if err != nil {
			res.Msg = "Bad Request Parsed Invalid CSRF"
			context.Abort()
			context.JSON(res.Code, res)
			return
		}

		if csrfToken.SessionId != session.ID {
			res.Code = http.StatusForbidden
			res.Msg = "Forbidden - Invalid CSRF"
			context.Abort()
			context.JSON(res.Code, res)
			return
		}

		if csrfToken.MaxAge.Before(time.Now()) {
			res.Code = http.StatusForbidden
			res.Msg = "Forbidden - Expired CSRF"
			context.Abort()
			context.JSON(res.Code, res)
			return
		}

		context.Next()

		//if session.Values != nil {
		//	isLogin := session.Values["login"]
		//	if isLogin == false || isLogin == nil {
		//		context.Next() // 如果 session 没有登陆则默认放行所有请求
		//		return
		//	}
		//} else {
		//	context.Next() // 如果 session 没有登陆则默认放行所有请求
		//	return
		//}

		//if csrf == "" { //测试时使用自动生成, 实际中应该用 websocket 向客户端传输后给模板
		//	//tokenObject, _ := json.Marshal(csrfToken)
		//	//token := cipher.EncryptedToBase64(string(tokenObject))
		//	//http.SetCookie(context.Writer, &http.Cookie{
		//	//	Name:    "csrf",
		//	//	Value:   token,
		//	//	Expires: time.Now().Add(24 * time.Hour),
		//	//})
		//	//res.Code = http.StatusForbidden
		//	//res.Msg = "Forbidden - Please refresh the page. "
		//	//context.Abort()
		//	//context.JSON(res.Code, res)
		//	//return
		//} else {
		//
		//}
	}
}
