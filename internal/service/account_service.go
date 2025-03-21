package service

import (
	"my_destributed_project/internal/models"
	"my_destributed_project/internal/repo"
)

type AccountService struct {
	Repo *repo.AccountRepo
}

// GetAccount 获取账户
func (s *AccountService) GetAccount(id uint) (*models.Account, error) {
	return s.Repo.GetAccountByID(id)
}

func (s *AccountService) CreateAccount(account *models.Account) error {
	return s.Repo.CreateAccount(account)
}

func (s *AccountService) EditAccount(account *models.Account) error {
	return s.Repo.EditAccount(account)
}

func (s *AccountService) DelAccount(id uint) error {
	return s.Repo.DelAccount(id)
}
