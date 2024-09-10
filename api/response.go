package api

import (
	"crawlerctl/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SuccessResponse 返回成功的 JSON 响应
func SuccessResponse(c *gin.Context, data interface{}) {
	response := models.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// ErrorResponse 返回错误的 JSON 响应
func ErrorResponse(c *gin.Context, status int, message string) {
	response := models.Response{
		Status:  status,
		Message: message,
		Data:    nil,
	}
	c.JSON(status, response)
}
