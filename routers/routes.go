package routes

import (
	"github.com/Ntrashh/crawlerctl/api"
	"github.com/Ntrashh/crawlerctl/middleware"
	"github.com/Ntrashh/crawlerctl/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(router *gin.Engine) {

	router.POST("/login", api.LoginHandler)
	// 环境管理相关的路由
	envService := services.NewEnvService()
	envHandler := api.NewEnvHandler(envService)
	envRoutes := router.Group("/envs")
	envRoutes.Use(middleware.AuthMiddleware())
	{
		envRoutes.GET("/check_install", envHandler.CheckPyenvInstalledHandler)
		envRoutes.GET("/get_versions", envHandler.GetPyenvPythonVersionHandler)
		envRoutes.POST("/install", envHandler.InstallPythonHandler)
		envRoutes.GET("/remote_versions", envHandler.GetRemotePythonVersionHandler)
		envRoutes.POST("/set_global", envHandler.SetVersionGlobalHandler)
		envRoutes.POST("/delete_python", envHandler.DeletePythonVersionHandler)
		envRoutes.POST("/create_virtualenv", envHandler.CreateVirtualenvHandler)
		envRoutes.POST("/delete_virtualenv", envHandler.DeleteVirtualenvHandler)
	}

	taskRoutes := router.Group("/tasks")
	{
		taskRoutes.GET("/task_status", api.GetTaskStatus)
	}

}
