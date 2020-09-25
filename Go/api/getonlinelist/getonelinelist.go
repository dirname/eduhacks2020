package getonlinelist

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
	GroupName string      `json:"groupName" validate:"required"`
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
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

	systemID := context.Request.Header.Get("SystemID")
	ret := websocket.GetOnlineList(&systemID, &inputData.GroupName)

	api.Render(context.Writer, retcode.SUCCESS, "success", ret)
	return
}
