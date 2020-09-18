package response

// TableResponse 表格渲染的数据结构
type TableResponse struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
	Count   int64       `json:"count"`
}
