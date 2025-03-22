package handlers

import (
	"github.com/gin-gonic/gin"
	"my_destributed_project/internal/dtos"
	"my_destributed_project/pkg/log"
	"my_destributed_project/pkg/utils"
	"net/http"
)

type HelloHandler struct {
}

// HelloHandler 处理 /hello 请求
func (h *HelloHandler) HelloHandler(c *gin.Context) {
	log.Logger.Info("Hello handler called")
	c.JSON(http.StatusOK, gin.H{"message": "hello world!"})
}

// TestMq 测试Kafka消息队列
func (h *HelloHandler) TestMq(c *gin.Context) {
	var dto dtos.MqDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	topic := dto.Topic
	message := dto.Message

	err := utils.ProduceMessage(topic, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "Message sent!"})
	}
}
