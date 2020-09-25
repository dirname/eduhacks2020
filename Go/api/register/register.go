package register

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
	SystemID string `json:"systemId" validate:"required"`
}

// Run 运行
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

	err = websocket.Register(inputData.SystemID)
	if err != nil {
		api.Render(context.Writer, retcode.FAIL, err.Error(), []string{})
		return
	}

	api.Render(context.Writer, retcode.SUCCESS, "success", []string{})
	return
}
