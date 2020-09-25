package send2clients

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/define/retcode"
	"eduhacks2020/Go/protocol/websocket"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Controller 控制器
type Controller struct {
}

type inputData struct {
	ClientIds  []string `json:"clientIds" validate:"required"`
	SendUserID string   `json:"sendUserId"`
	Code       int      `json:"code"`
	Msg        string   `json:"msg"`
	Data       string   `json:"data"`
}

// Run 运行实例
func (c *Controller) Run(context *gin.Context) {
	var inputData inputData
	if err := json.NewDecoder(context.Request.Body).Decode(&inputData); err != nil {
		context.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err := api.Validate(inputData)
	if err != nil {
		api.Render(context.Writer, retcode.FAIL, err.Error(), []string{})
		return
	}

	for _, clientID := range inputData.ClientIds {
		//发送信息
		_ = websocket.SendMessage2Client(clientID, inputData.SendUserID, inputData.Code, inputData.Msg, &inputData.Data)
	}

	api.Render(context.Writer, retcode.SUCCESS, "success", []string{})
	return
}
