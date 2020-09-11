package models

import "time"

//普通 http 请求的响应体
type Response struct {
	Code   int         `json:"code"`
	Path   string      `json:"path"`
	Method string      `json:"method"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
	Time   time.Time   `json:"time"`
	IP     string      `json:"ip"`
}
