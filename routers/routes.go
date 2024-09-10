package routes

import (
	"crawlerctl/api"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {
	// 环境管理相关的路由
	envRoutes := router.Group("/envs")
	{
		envRoutes.GET("/check_install", api.CheckPyenvInstalledHandler)
		envRoutes.GET("/get_versions", api.GetPyenvPythonVersionHandler)
	}

	// 其他路由
	// ...
}
