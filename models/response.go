package models

// Response 是 API 统一的返回结构
type Response struct {
	Status  int         `json:"status"`  // HTTP 状态码
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 数据
}
