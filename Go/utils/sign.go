package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
)

// 定义一些常量错误
const (
	SignInvalid = "data has been tampered with invalid sign" // 签名不正确的错误信息
)

// 用于计算签名
func calcSign(timestamp string, data []byte) string {
	var buffer bytes.Buffer
	buffer.Write([]byte(timestamp))
	buffer.Write(data)
	h := md5.New()
	h.Write(buffer.Bytes())
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// VerifySign 返回校验的结果
func VerifySign(salt string, submitSign, data []byte) bool {
	submit := string(submitSign)
	calc := calcSign(salt, data)
	return submit == calc
}
