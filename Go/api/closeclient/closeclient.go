package closeclient

import (
	"eduhacks2020/Go/api"
	"eduhacks2020/Go/define/retcode"
	"eduhacks2020/Go/protocol/websocket"
	"encoding/json"
	"net/http"
)

// Controller
type Controller struct {
}

type inputData struct {
	ClientID string `json:"clientId" validate:"required"`
}

// Run 启动路由
func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	var inputData inputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := api.Validate(inputData)
	if err != nil {
		api.Render(w, retcode.FAIL, err.Error(), []string{})
		return
	}

	systemID := r.Header.Get("SystemID")

	//发送信息
	websocket.CloseClient(inputData.ClientID, systemID)

	api.Render(w, retcode.SUCCESS, "success", map[string]string{})
	return
}
