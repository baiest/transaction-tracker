package usecase

import (
	"context"
	"transaction-tracker/internal/accounts/domain"

	"github.com/golang-jwt/jwt/v5"
)

type AccountsUseCase interface {
	GetAuthURL() string
	CreateAccount(ctx context.Context, account *domain.Account) error
	GetAccount(ctx context.Context, accountID string) (*domain.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	GetOrCreateAccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	SaveGoogleAccount(ctx context.Context, accountID string, code string) error
	GenerateTokens(ctx context.Context, account *domain.Account) (string, string, string, error)
	CreateWatcher(ctx context.Context, account *domain.Account) error
	DeleteWatcher(ctx context.Context, account *domain.Account) error
	VerifyToken(tokenString string) (*jwt.Token, error)
}
