package routes

import (
	"crawlerctl/api"
	"crawlerctl/middleware"
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
	}

}
