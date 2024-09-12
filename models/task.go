package models

type Task struct {
	ID      string                                        `json:"id"`
	Params  interface{}                                   `json:"params"` // 任务参数
	Status  string                                        `json:"status"` // 任务状态
	Result  interface{}                                   `json:"result"` // 任务结果
	Info    interface{}                                   `json:"info"`
	Err     error                                         `json:"-"` // 错误
	Execute func(params interface{}) (interface{}, error) // 执行逻辑
}
