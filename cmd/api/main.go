package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"my_destributed_project/configs"
	"my_destributed_project/internal/app/handlers"
	"my_destributed_project/pkg/database"
	"my_destributed_project/pkg/log"
)

func main() {
	// 初始化日志
	log.InitLogger()
	defer log.Logger.Sync()

	// 加载配置
	if err := configs.LoadConfig("configs/config.yaml"); err != nil {
		log.Logger.Fatal("Failed to load config", zap.Error(err))
	}

	// 连接数据库
	if err := database.ConnectDatabase(); err != nil {
		log.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 初始化 Gin 路由
	r := gin.Default()

	// 注册路由
	r.GET("/hello", handlers.HelloHandler)
	r.GET("/user/:id", handlers.UserHandler)

	// 启动服务器
	port := ":8080"
	fmt.Println("Server running on", port)
	if err := r.Run(port); err != nil {
		log.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
