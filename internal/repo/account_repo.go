package repo

import (
	"gorm.io/gorm"
	"my_destributed_project/internal/models"
)

// AccountRepo 负责数据库操作
type AccountRepo struct {
	DB *gorm.DB
}

// GetAccountByID 根据 ID 查询账户
func (r *AccountRepo) GetAccountByID(id uint) (*models.Account, error) {
	var account models.Account
	if err := r.DB.First(&account, id).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepo) CreateAccount(account *models.Account) error {
	return r.DB.Create(account).Error
}

func (r *AccountRepo) EditAccount(account *models.Account) error {
	//sql := "UPDATE accounts SET account_holder_name = ?, balance = ? WHERE account_id = ?"
	//return r.DB.Exec(sql, account.AccountHolderName, account.Balance, account.AccountID).Error
	return r.DB.Save(account).Error
}

func (r *AccountRepo) DelAccount(id uint) error {
	// return r.DB.Where("account_number = ?", accountNumber).Delete(&models.Account{}).Error
	return r.DB.Delete(&models.Account{}, id).Error
}
