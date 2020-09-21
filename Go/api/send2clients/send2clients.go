package send2clients

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
	ClientIds  []string `json:"clientIds" validate:"required"`
	SendUserID string   `json:"sendUserId"`
	Code       int      `json:"code"`
	Msg        string   `json:"msg"`
	Data       string   `json:"data"`
}

// Run
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

	for _, clientID := range inputData.ClientIds {
		//发送信息
		_ = websocket.SendMessage2Client(clientID, inputData.SendUserID, inputData.Code, inputData.Msg, &inputData.Data)
	}

	api.Render(w, retcode.SUCCESS, "success", []string{})
	return
}
