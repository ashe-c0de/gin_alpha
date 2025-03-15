package handlers

import (
	"github.com/gin-gonic/gin"
	"my_destributed_project/pkg/log"
	"net/http"
)

// HelloHandler 处理 /hello 请求
func HelloHandler(c *gin.Context) {
	log.Logger.Info("Hello handler called")
	c.JSON(http.StatusOK, gin.H{"message": "hello world!"})
}
