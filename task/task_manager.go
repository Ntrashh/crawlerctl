package task

import (
	"github.com/Ntrashh/crawlerctl/models"
	"github.com/google/uuid"
	"sync"
)

var TaskStore sync.Map // 全局任务存储

// 启动任务并传入执行逻辑和参数
func StartTask(execute func(params interface{}) (interface{}, error), params interface{}) string {
	taskID := uuid.New().String()

	// 创建并初始化任务
	task := &models.Task{
		ID:      taskID,
		Params:  params,
		Status:  "pending",
		Execute: execute,
	}

	TaskStore.Store(taskID, task)

	// 后台执行任务
	go runTask(task)

	return taskID
}

// 执行任务的函数
func runTask(task *models.Task) {
	task.Status = "running"
	TaskStore.Store(task.ID, task)

	// 执行传入的任务逻辑
	result, err := task.Execute(task.Params)

	TaskStore.Store(task.ID, task) // 更新任务状态和结果
	// 根据执行结果更新任务状态
	if err != nil {
		task.Status = "failed"
		task.Err = err
	} else {
		task.Status = "done"
	}
	task.Result = result
	TaskStore.Store(task.ID, task)
}

// 根据任务ID获取任务状态
func GetTaskStatus(taskID string) (*models.Task, bool) {
	if task, ok := TaskStore.Load(taskID); ok {
		return task.(*models.Task), true
	}
	return nil, false
}
