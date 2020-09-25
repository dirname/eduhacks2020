package closeclient

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
	ClientID string `json:"clientId" validate:"required"`
}

// Run 启动路由
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

	systemID := context.Request.Header.Get("SystemID")

	//发送信息
	websocket.CloseClient(inputData.ClientID, systemID)

	api.Render(context.Writer, retcode.SUCCESS, "success", map[string]string{})
	return
}
