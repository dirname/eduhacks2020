package protocol

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	//"go-websocket/api"
	//"go-websocket/define/retcode"
	//"go-websocket/tools/util"
	"net/http"
)

const (
	// 最大的消息大小
	maxMessageSize = 8192
)

type Controller struct {
}

type renderData struct {
	ClientId string `json:"clientId"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("upgrade error: %v", err)
		http.NotFound(w, r)
		return
	}

	//设置读取消息大小上线
	conn.SetReadLimit(maxMessageSize)

	//解析参数
	//systemId := r.FormValue("systemId")
	//if len(systemId) == 0 {
	//	_ = Render(conn, "", "", retcode.SYSTEM_ID_ERROR, "系统ID不能为空", []string{})
	//	_ = conn.Close()
	//	return
	//}
	//
	//clientId := util.GenClientId()

	clientSocket := NewClient("", "", conn)

	//Manager.AddClient2SystemClient(systemId, clientSocket)

	//读取客户端消息
	clientSocket.Read()

	//test := &protobuf.Student{
	//	Name:    "<link rel=\"stylesheet\" href=\"../layui/css/layui.css\">\n<body>\n \n<!-- 你的HTML代码 -->\n \n<script src=\"../layui/layui.all.js\"></script>\n<script>\n//一般直接写在一个js文件中\nlayui.use(['layer', 'form'], function(){\n  var layer = layui.layer\n  ,form = layui.form;\n  \n  layer.msg('Hello World');\n});\n</script>",
	//	Comment: []byte("Data"),
	//	Scores:  []int32{98, 85, 88},
	//}
	//data, err := proto.Marshal(test)
	//conn.WriteMessage(2, data)

	//if err = api.ConnRender(conn, renderData{ClientId: clientId}); err != nil {
	//	_ = conn.Close()
	//	return
	//}
	//
	//// 用户连接事件
	//Manager.Connect <- clientSocket
}
