package bind2group

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
	ClientID  string `json:"clientId" validate:"required"`
	GroupName string `json:"groupName" validate:"required"`
	UserID    string `json:"userId"`
	Extend    string `json:"extend"` // 拓展字段，方便业务存储数据
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
	websocket.AddClient2Group(systemID, inputData.GroupName, inputData.ClientID, inputData.UserID, inputData.Extend)

	api.Render(context.Writer, retcode.SUCCESS, "success", []string{})
}
