package handlers

import (
	"github.com/gin-gonic/gin"
	"my_destributed_project/internal/app/service"
	"my_destributed_project/pkg/log"
	"net/http"
	"strconv"
)

// UserHandler 处理 /user/:id 请求
func UserHandler(c *gin.Context) {
	log.Logger.Info("account handler called")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := service.GetAccount(uint(id))
	if err != nil {
		log.Logger.Error("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
