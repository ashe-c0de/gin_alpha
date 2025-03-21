package service

import (
	"gorm.io/gorm"
	"my_destributed_project/internal/models"
	"my_destributed_project/internal/repo"
)

type AccountService struct {
	Repo *repo.AccountRepo
	DB   *gorm.DB // 需要传入 GORM 的 DB 实例，用于事务管理
}

// GetAccount 获取账户（无需事务）
func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	return s.Repo.GetAccountByID(id)
}

// CreateAccount 在事务中创建账户
func (s *AccountService) CreateAccount(account *models.Account) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 使用事务对象 tx 代替 r.DB
		if err := s.Repo.CreateAccountTx(tx, account); err != nil {
			return err // 发生错误时 GORM 会自动回滚事务
		}

		// 如果有额外的业务逻辑（例如记录日志），可以在这里添加
		// log := models.AccountLog{AccountID: account.ID, Action: "created"}
		// if err := s.Repo.LogAccountActionTx(tx, &log); err != nil {
		// 	return err
		// }

		return nil // 事务提交
	})
}

// EditAccount 在事务中更新账户
func (s *AccountService) EditAccount(account *models.Account) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		// 使用事务对象 tx 代替 r.DB
		if err := s.Repo.EditAccountTx(tx, account); err != nil {
			return err
		}
		return nil
	})
}

// DelAccount 删除账户（可以加事务，但通常 DELETE 操作不需要事务）
// 如果删除操作涉及 多个表的级联删除，或者 有额外的业务逻辑，建议使用事务
func (s *AccountService) DelAccount(id uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := s.Repo.DelAccountTx(tx, id); err != nil {
			return err
		}
		return nil
	})
}
