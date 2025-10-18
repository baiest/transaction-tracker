package repository

import (
	"context"
	"testing"

	"transaction-tracker/internal/accounts/domain"
	database "transaction-tracker/pkg/databases/mongo"
	"transaction-tracker/pkg/google"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestAccountsRepository_CreateAccount(t *testing.T) {
	c := require.New(t)

	mockCollection := &database.MockCollection{
		UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
			return &mongo.UpdateResult{}, nil
		},
	}

	repo := NewAccountsRepository(mockCollection)

	account, err := domain.NewAccount("test@example.com")
	c.NoError(err)

	err = repo.CreateAccount(context.Background(), account)

	c.NoError(err)
}

func TestAccountsRepository_GetAccount(t *testing.T) {
	c := require.New(t)

	account := &domain.Account{ID: "test-id"}

	t.Run("success", func(t *testing.T) {
		mockCollection := &database.MockCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(account, nil, nil)
			},
		}

		repo := NewAccountsRepository(mockCollection)

		result, err := repo.GetAccount(context.Background(), account.ID)
		c.NoError(err)
		c.Equal(account, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockCollection := &database.MockCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return &mongo.SingleResult{}
			},
		}

		repo := NewAccountsRepository(mockCollection)

		result, err := repo.GetAccount(context.Background(), account.ID)
		c.Error(err)
		c.Equal(ErrAccountNotfound, err)
		c.Nil(result)
	})
}

func TestAccountsRepository_GetAccountByEmail(t *testing.T) {
	c := require.New(t)

	account := &domain.Account{Email: "test@example.com"}

	t.Run("success", func(t *testing.T) {
		mockCollection := &database.MockCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return mongo.NewSingleResultFromDocument(account, nil, nil)
			},
		}

		repo := NewAccountsRepository(mockCollection)

		result, err := repo.GetAccountByEmail(context.Background(), account.Email)

		c.NoError(err)
		c.Equal(account, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockCollection := &database.MockCollection{
			FindOneFn: func(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult {
				return &mongo.SingleResult{}
			},
		}

		repo := NewAccountsRepository(mockCollection)

		result, err := repo.GetAccountByEmail(context.Background(), account.Email)

		c.Error(err)
		c.Equal(ErrAccountNotfound, err)
		c.Nil(result)
	})
}

func TestAccountsRepository_SaveGoogleAccount(t *testing.T) {
	c := require.New(t)

	mockCollection := &database.MockCollection{
		UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
			return &mongo.UpdateResult{}, nil
		},
	}

	repo := NewAccountsRepository(mockCollection)

	accountID := "test-id"
	googleAccount := &google.GoogleAccount{}

	err := repo.SaveGoogleAccount(context.Background(), accountID, googleAccount)
	c.NoError(err)
}

func TestAccountsRepository_UpdateAccount(t *testing.T) {
	c := require.New(t)

	mockCollection := &database.MockCollection{
		UpdateOneFn: func(ctx context.Context, filter, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error) {
			return &mongo.UpdateResult{}, nil
		},
	}

	repo := NewAccountsRepository(mockCollection)

	account := &domain.Account{ID: "test-id"}

	err := repo.UpdateAccount(context.Background(), account)

	c.NoError(err)
}
