package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"my_destributed_project/api"
	"my_destributed_project/configs"
	"my_destributed_project/internal/handlers"
	"my_destributed_project/internal/repo"
	"my_destributed_project/internal/service"
	"my_destributed_project/pkg/database"
	"my_destributed_project/pkg/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("start")

	// 加载配置
	if err := configs.LoadConfig("configs/config.yaml"); err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}
	fmt.Println("Success to load config")

	// 初始化日志
	log.InitLogger()
	defer func() {
		if err := log.Logger.Sync(); err != nil {
			fmt.Println("Failed to sync logger:", err)
		}
	}()

	// 连接数据库
	db := database.ConnectDatabase()

	// DI account handler
	accountRepo := &repo.AccountRepo{DB: db}
	accountService := &service.AccountService{Repo: accountRepo}
	accountHandler := &handlers.AccountHandler{Service: accountService}

	// DI hello handler
	helloHandler := &handlers.HelloHandler{}

	// 设置 API 路由
	routers := api.SetRouters(helloHandler, accountHandler)

	port := configs.AppConfig.Server.Port

	server := &http.Server{
		Addr:    ":" + port,
		Handler: routers,
	}

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动 HTTP 服务器
	ctx, srvCancel := context.WithCancel(context.Background())

	go func() {
		log.Logger.Info("Server is starting: ", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Error("Server startup failed", zap.Error(err))
			srvCancel() // 触发 `context` 取消，确保 `server.Shutdown(ctx)` 被调用
		}
	}()

	// 等待退出信号
	<-quit
	log.Logger.Info("Shutting down server...")

	// 关闭 HTTP 服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Logger.Fatal("Server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	log.Logger.Info("Server exited properly")
}
