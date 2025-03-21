package dtos

import "my_destributed_project/internal/models"

// AccountDTO 用于 API 交互的数据结构
type AccountDTO struct {
	AccountID         uint    `json:"account_id"`                             // 账户 ID
	AccountNumber     string  `json:"account_number" binding:"required"`      // 账户号（必填）
	AccountHolderName string  `json:"account_holder_name" binding:"required"` // 账户持有人姓名（必填）
	Balance           float64 `json:"balance" binding:"required,gte=0"`       // 账户余额（必填，且不能小于 0）
}

// ToAccountEntity 将 DTO 转换为数据库模型（用于保存到数据库）
func ToAccountEntity(dto AccountDTO) models.Account {
	return models.Account{
		AccountNumber:     dto.AccountNumber,
		AccountHolderName: dto.AccountHolderName,
		Balance:           dto.Balance,
	}
}

type DelAccountDTO struct {
	AccountID uint `json:"account_id"` // 账户 ID
}
