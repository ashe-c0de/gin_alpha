package main

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"my_destributed_project/api"
	"my_destributed_project/configs"
	"my_destributed_project/internal/etcd"
	"my_destributed_project/internal/handlers"
	"my_destributed_project/internal/repo"
	"my_destributed_project/internal/service"
	"my_destributed_project/pkg/database"
	"my_destributed_project/pkg/log"
	"my_destributed_project/pkg/utils"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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
	accountService := &service.AccountService{Repo: accountRepo, DB: db}
	accountHandler := &handlers.AccountHandler{Service: accountService}

	// DI hello handler
	helloHandler := &handlers.HelloHandler{}

	// 设置 API 路由
	routers := api.SetRouters(helloHandler, accountHandler)

	port := configs.AppConfig.Server.Port

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: routers,
	}

	// 监听退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 创建一个 ctx，用于控制 HTTP 服务器的生命周期
	ctx, srvCancel := context.WithCancel(context.Background())

	// 注册服务
	etcd.Init()
	etcd.RegisterService("/services/my-app/instance-1", "http://139.196.243.6:8000", 10)

	// 启动 HTTP 服务器
	go func() {
		log.Logger.Info("Server is starting", zap.Int("port", port))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Logger.Error("Server startup failed", zap.Error(err))
			srvCancel() // 触发 `context` 取消，确保 `server.Shutdown(ctx)` 被调用
		}
	}()

	// 启动Kafka消费者
	go utils.StartConsumer(ctx, utils.GroupID)

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
