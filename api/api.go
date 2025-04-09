package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"my_destributed_project/internal/handlers"
)

func SetRouters(helloHandler *handlers.HelloHandler, accountHandler *handlers.AccountHandler) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	// 初始化 Gin 路由
	routers := gin.New()

	// 跨域配置
	routers.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// 注册路由
	routers.GET("/hello", helloHandler.HelloHandler)
	routers.POST("/mq", helloHandler.TestMq)

	routers.GET("/acc/:id", accountHandler.GetAccount)
	routers.POST("/acc/add", accountHandler.CreatAccount)
	routers.POST("/acc/edit", accountHandler.EditAccount)
	routers.POST("/acc/del", accountHandler.DelAccount)

	return routers
}
