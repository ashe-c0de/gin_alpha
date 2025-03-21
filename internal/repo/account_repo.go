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

// CreateAccountTx 使用事务创建账户
func (r *AccountRepo) CreateAccountTx(tx *gorm.DB, account *models.Account) error {
	return tx.Create(account).Error
}

// EditAccountTx 使用事务更新账户
func (r *AccountRepo) EditAccountTx(tx *gorm.DB, account *models.Account) error {
	//sql := "UPDATE accounts SET account_holder_name = ?, balance = ? WHERE account_id = ?"
	//return tx.Exec(sql, account.AccountHolderName, account.Balance, account.AccountID).Error
	return tx.Save(account).Error
}

// DelAccountTx 使用事务删除账户
func (r *AccountRepo) DelAccountTx(tx *gorm.DB, id uint) error {
	//return tx.Where("account_number = ?", id).Delete(&models.Account{}).Error
	return tx.Delete(&models.Account{}, id).Error
}
