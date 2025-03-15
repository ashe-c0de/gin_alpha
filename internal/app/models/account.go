package models

import (
	"time"
)

// Account 账户表实体
type Account struct {
	AccountID         uint      `gorm:"primaryKey;column:account_id"` // 账户 ID
	AccountNumber     string    `gorm:"unique;column:account_number"` // 账户号（唯一）
	AccountHolderName string    `gorm:"column:account_holder_name"`   // 账户持有人姓名
	Balance           float64   `gorm:"column:balance"`               // 账户余额
	CreatedAt         time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定数据库表名
func (Account) TableName() string {
	return "account"
}
