package handlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"my_destributed_project/internal/dtos"
	"my_destributed_project/internal/service"
	"my_destributed_project/pkg/log"
	"net/http"
	"strconv"
)

type AccountHandler struct {
	Service *service.AccountService
}

// GetAccount 处理 /user/:id 请求
func (h *AccountHandler) GetAccount(c *gin.Context) {
	log.Logger.Info("account handler called")
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Logger.Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.Service.GetAccount(uint(id))
	if err != nil {
		log.Logger.Error("User not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreatAccount 创建账户
func (h *AccountHandler) CreatAccount(c *gin.Context) {
	log.Logger.Info("start create account")
	var dto dtos.AccountDTO

	// 校验 & 转换为数据库实体
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	account := dtos.ToAccountEntity(dto)

	log.Logger.Info("创建账户", zap.Any("Account", account))

	// 新增Account
	err := h.Service.CreateAccount(&account)
	if err != nil {
		log.Logger.Error("创建account失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "success",
		})
	}
}

func (h *AccountHandler) EditAccount(c *gin.Context) {
	var dto dtos.AccountDTO

	// 校验 & 转换为数据库实体
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 查询是否存在该账户
	account, _ := h.Service.GetAccount(dto.AccountID)
	if account == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "修改的账户不存在"})
		return
	}
	account.AccountNumber = dto.AccountNumber
	account.AccountHolderName = dto.AccountHolderName
	account.Balance = dto.Balance
	// 新增Account
	err := h.Service.EditAccount(account)
	if err != nil {
		log.Logger.Error("修改account失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "success",
		})
	}
}

func (h *AccountHandler) DelAccount(c *gin.Context) {
	var dto dtos.DelAccountDTO
	// 校验 & 转换为数据库实体
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.Service.DelAccount(dto.AccountID)
	if err != nil {
		log.Logger.Error("删除account失败", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": "200",
			"msg":  "del success",
		})
	}
}
