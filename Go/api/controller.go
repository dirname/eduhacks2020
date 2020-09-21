package api

import (
	"eduhacks2020/Go/define/retcode"
	"eduhacks2020/Go/protobuf"
	"encoding/base64"
	"encoding/json"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	zh2 "gopkg.in/go-playground/validator.v9/translations/zh"
	"io"
	"net/http"
)

type RetData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func xorData(data []byte, decrypt bool) []byte {
	res := make([]byte, len(data))
	for i, b := range data {
		res[i] = b ^ 32
	}
	if !decrypt {
		return []byte(base64.URLEncoding.EncodeToString(res))
	}
	return res
}

func ConnRender(conn *websocket.Conn, data interface{}) (err error) {
	js, _ := json.Marshal(RetData{
		Code: retcode.SUCCESS,
		Msg:  "success",
		Data: data,
	})
	res := &protobuf.Response{
		Code:   0,
		Msg:    "OK",
		Type:   1,
		Data:   js,
		Render: false,
		Html:   nil,
		Id:     "",
	}
	msg, _ := proto.Marshal(res)
	err = conn.WriteMessage(2, xorData(msg, false))
	return
}

func Render(w http.ResponseWriter, code int, msg string, data interface{}) (str string) {
	var retData RetData

	retData.Code = code
	retData.Msg = msg
	retData.Data = data

	retJson, _ := json.Marshal(retData)
	str = string(retJson)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = io.WriteString(w, str)
	return
}

func Validate(inputData interface{}) error {

	validate := validator.New()
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")

	_ = zh2.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(inputData)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(trans))
		}
	}

	return nil
}
