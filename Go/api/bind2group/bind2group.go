package bind2group

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
	ClientID  string `json:"clientId" validate:"required"`
	GroupName string `json:"groupName" validate:"required"`
	UserID    string `json:"userId"`
	Extend    string `json:"extend"` // 拓展字段，方便业务存储数据
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
	websocket.AddClient2Group(systemID, inputData.GroupName, inputData.ClientID, inputData.UserID, inputData.Extend)

	api.Render(w, retcode.SUCCESS, "success", []string{})
}
