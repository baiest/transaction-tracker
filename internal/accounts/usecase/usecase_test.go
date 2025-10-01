package usecase

import (
	"context"
	"testing"

	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/google"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockAccountsRepository is a mock of AccountsRepository
type MockAccountsRepository struct {
	mock.Mock
}

func (m *MockAccountsRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *MockAccountsRepository) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountsRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountsRepository) SaveGoogleAccount(ctx context.Context, accountID string, googleAccount *google.GoogleAccount) error {
	args := m.Called(ctx, accountID, googleAccount)
	return args.Error(0)
}

func (m *MockAccountsRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func TestNewAccountsUseCase(t *testing.T) {
	c := require.New(t)

	googleClient := &google.GoogleClient{}
	repo := new(MockAccountsRepository)

	uc := NewAccountsUseCase(googleClient, repo)

	c.NotNil(uc)
}

func TestAccountsUseCase_GetAuthURL(t *testing.T) {
	c := require.New(t)

	googleClient := &google.MockGoogleClient{}
	repo := new(MockAccountsRepository)
	uc := NewAccountsUseCase(googleClient, repo)

	googleClient.On("GetAuthURL").Return("http://test.com")

	url := uc.GetAuthURL()

	c.Equal("http://test.com", url)
	googleClient.AssertExpectations(t)
}

func TestAccountsUseCase_CreateAccount(t *testing.T) {
	c := require.New(t)

	googleClient := &google.GoogleClient{}
	repo := new(MockAccountsRepository)
	uc := NewAccountsUseCase(googleClient, repo)

	account := &domain.Account{}
	repo.On("CreateAccount", context.Background(), account).Return(nil)

	err := uc.CreateAccount(context.Background(), account)

	c.NoError(err)
	repo.AssertExpectations(t)
}

func TestAccountsUseCase_GetAccount(t *testing.T) {
	c := require.New(t)

	googleClient := &google.GoogleClient{}
	repo := new(MockAccountsRepository)
	uc := NewAccountsUseCase(googleClient, repo)

	account := &domain.Account{ID: "test-id"}
	repo.On("GetAccount", context.Background(), account.ID).Return(account, nil)

	result, err := uc.GetAccount(context.Background(), account.ID)

	c.NoError(err)
	c.Equal(account, result)
	repo.AssertExpectations(t)
}

func TestAccountsUseCase_GetAccountByEmail(t *testing.T) {
	c := require.New(t)

	googleClient := &google.GoogleClient{}
	repo := new(MockAccountsRepository)
	uc := NewAccountsUseCase(googleClient, repo)

	account := &domain.Account{Email: "test@example.com"}

	t.Run("success", func(t *testing.T) {
		repo.On("GetAccountByEmail", context.Background(), account.Email).Return(account, nil).Once()

		result, err := uc.GetAccountByEmail(context.Background(), account.Email)

		c.NoError(err)
		c.Equal(account, result)
		repo.AssertExpectations(t)
	})

	t.Run("empty email", func(t *testing.T) {
		result, err := uc.GetAccountByEmail(context.Background(), "")

		c.Error(err)
		c.Nil(result)
		c.Equal("email cannot be empty", err.Error())
	})
}
