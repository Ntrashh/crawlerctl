package routes

import (
	"github.com/Ntrashh/crawlerctl/api"
	"github.com/Ntrashh/crawlerctl/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {

	router.POST("/login", api.LoginHandler)
	// 环境管理相关的路由
	envRoutes := router.Group("/envs")
	envRoutes.Use(middleware.AuthMiddleware())
	{
		envRoutes.GET("/check_install", api.CheckPyenvInstalledHandler)
		envRoutes.GET("/get_versions", api.GetPyenvPythonVersionHandler)
		envRoutes.POST("/install", api.InstallPythonHandler)
		envRoutes.GET("/remote_versions", api.GetRemotePythonVersionHandler)
		envRoutes.POST("/set_global", api.SetVersionGlobalHandler)
		envRoutes.POST("/delete_python", api.DeletePythonVersionHandler)
		envRoutes.POST("/create_virtualenv", api.CreateVirtualenvHandler)
		envRoutes.POST("/delete_virtualenv", api.DeleteVirtualenvHandler)
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("/task_status", api.GetTaskStatus)
	}

}
