package service

import (
	"my_destributed_project/internal/app/models"
	"my_destributed_project/internal/app/repo"
)

// GetAccount 获取账户
func GetAccount(id uint) (*models.Account, error) {
	return repo.GetAccountByID(id)
}
