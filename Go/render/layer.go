package render

import (
	"fmt"
	"io/ioutil"
	"os"
)

// 定义 Layer 的类型
const (
	None      = "undefined" //没有图标
	Alert     = 0           //一个❗的图标
	Right     = 1           //一个✔的图标
	Incorrect = 2           //一个❌的图标
	Doubt     = 3           //一个❓的图标
	Lock      = 4           //一个🔒的图标
	Sad       = 5           //一个😟的图标
	Smile     = 6           //一个😊的图标

)

// ReadTemp 读取模板
func ReadTemp(filePth string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filePth)
	fileString := string(fileBytes)
	return fileString, err
}

// GetLayer 获取一个 Layer 的代码
func GetLayer(t, icon int, title, content string) string {
	dir, _ := os.Getwd()
	temp, _ := ReadTemp(dir + "/render/html/layer.temp")
	return fmt.Sprintf(temp, t, title, content, icon)
}

// GetMsg 获取一个 MSg 的代码
func GetMsg(content string, Sec int) string {
	dir, _ := os.Getwd()
	temp, _ := ReadTemp(dir + "/render/html/msg.temp")
	return fmt.Sprintf(temp, content, Sec*1000)
}
