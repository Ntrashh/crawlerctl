package main

import (
	"github.com/Ntrashh/crawlerctl/config"
	"github.com/Ntrashh/crawlerctl/database"
	routes "github.com/Ntrashh/crawlerctl/routers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	config.LoadConfig()
	database.InitDatabase()

	// 初始化 Gin
	router := gin.Default()

	// 添加 CORS 中间件
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                // 允许的前端地址
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "DELETE"}, // 允许的请求方法
		AllowHeaders:     []string{"*"},                                // 允许的请求头
		AllowCredentials: true,                                         // 是否允许传递认证信息，比如 Cookies
	}))

	// 注册所有路由
	routes.RegisterRoutes(router)
	// 启动服务器
	log.Println("Server starting on port 8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
