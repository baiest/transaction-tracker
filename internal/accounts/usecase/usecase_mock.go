package usecase

import (
	"context"

	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/google"
)

// MockAccountsRepository es un mock del repository.AccountsRepository usado en usecase.
type MockAccountsRepository struct {
	CreateAccountFn      func(ctx context.Context, account *domain.Account) error
	GetAccountFn         func(ctx context.Context, id string) (*domain.Account, error)
	GetAccountByEmailFn  func(ctx context.Context, email string) (*domain.Account, error)
	SaveGoogleAccountFn  func(ctx context.Context, accountID string, ga *google.GoogleAccount) error
	UpdateAccountFn      func(ctx context.Context, account *domain.Account) error
	DeleteWatcherFn      func(ctx context.Context, accountID string) error
	CreateWatcherFn      func(ctx context.Context, accountID string) error
	RefreshGoogleTokenFn func(ctx context.Context, accountID string) error
}

func NewMockAccountsRepository() *MockAccountsRepository {
	return &MockAccountsRepository{}
}

func (m *MockAccountsRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	if m.CreateAccountFn != nil {
		return m.CreateAccountFn(ctx, account)
	}

	return nil
}

func (m *MockAccountsRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	if m.GetAccountFn != nil {
		return m.GetAccountFn(ctx, id)
	}

	return nil, nil
}

func (m *MockAccountsRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	if m.GetAccountByEmailFn != nil {
		return m.GetAccountByEmailFn(ctx, email)
	}

	return nil, nil
}

func (m *MockAccountsRepository) SaveGoogleAccount(ctx context.Context, accountID string, ga *google.GoogleAccount) error {
	if m.SaveGoogleAccountFn != nil {
		return m.SaveGoogleAccountFn(ctx, accountID, ga)
	}

	return nil
}

func (m *MockAccountsRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	if m.UpdateAccountFn != nil {
		return m.UpdateAccountFn(ctx, account)
	}

	return nil
}

func (m *MockAccountsRepository) DeleteWatcher(ctx context.Context, accountID string) error {
	if m.DeleteWatcherFn != nil {
		return m.DeleteWatcherFn(ctx, accountID)
	}

	return nil
}

func (m *MockAccountsRepository) CreateWatcher(ctx context.Context, accountID string) error {
	if m.CreateWatcherFn != nil {
		return m.CreateWatcherFn(ctx, accountID)
	}

	return nil
}

func (m *MockAccountsRepository) RefreshGoogleToken(ctx context.Context, accountID string) error {
	if m.RefreshGoogleTokenFn != nil {
		return m.RefreshGoogleTokenFn(ctx, accountID)
	}

	return nil
}
