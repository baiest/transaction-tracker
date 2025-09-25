package repository

import (
	"context"
	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/google"
)

// AccountsRepository defines the operations to create and query accounts.
type AccountsRepository interface {
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, accountID string) (*domain.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	SaveGoogleAccount(ctx context.Context, accountID string, googleAccount *google.GoogleAccount) error
}
