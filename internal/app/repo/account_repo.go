package repo

import (
	"my_destributed_project/internal/app/models"
	"my_destributed_project/pkg/database"
)

// GetAccountByID 根据 ID 查询账户
func GetAccountByID(id uint) (*models.Account, error) {
	var account models.Account
	if err := database.DB.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}
