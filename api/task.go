package api

import (
	crawler "github.com/Ntrashh/crawlerctl/task"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTaskStatus(c *gin.Context) {
	taskID := c.Query("task_id")
	task, found := crawler.GetTaskStatus(taskID)
	if !found {
		ErrorResponse(c, http.StatusNotFound, "Task not found")
		return
	}
	// 返回任务状态和结果
	SuccessResponse(c, gin.H{
		"task_id": task.ID,
		"status":  task.Status,
		"result":  task.Result,
		"error":   task.Err,
	})
}
