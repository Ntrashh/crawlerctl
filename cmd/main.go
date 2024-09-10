package main

import (
	routes "crawlerctl/routers"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 初始化 Gin
	router := gin.Default()

	// 注册所有路由
	routes.RegisterRoutes(router)
	// 启动服务器
	log.Println("Server starting on port 8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
