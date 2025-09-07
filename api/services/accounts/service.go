package accounts

import (
	"context"
	"errors"
	"strings"
	"transaction-tracker/api/repositories"
	"transaction-tracker/database/mongo/schemas"
)

type AccountService struct {
	repo repositories.IAccountsRepository
}

var (
	ErrAccountNotFound = errors.New("account not found")
)

func NewAccountService(ctx context.Context) (*AccountService, error) {
	accountRepo, err := repositories.NewAccountsRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &AccountService{
		repo: accountRepo,
	}, nil
}

func (a *AccountService) CreateAccount(ctx context.Context, account *Account) error {
	acc := &schemas.Account{
		ID:           account.ID,
		Email:        account.Email,
		RefreshToken: account.RefreshToken,
	}

	return a.repo.CreateAccount(ctx, acc)
}

func (a *AccountService) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	account, err := a.repo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &Account{
		ID:    account.ID,
		Email: account.Email,
	}, nil
}

func (a *AccountService) GetAccountByEmail(ctx context.Context, email string) (*Account, error) {
	account, err := a.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "no documents in result") {
			return nil, ErrAccountNotFound
		}

		return nil, err
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	return &Account{
		ID:    account.ID,
		Email: account.Email,
	}, nil
}
