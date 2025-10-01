package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/google"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockCollection is a mock of CollectionAPI
type MockCollection struct {
	mock.Mock
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	args := m.Called(ctx, filter, update, opts)
	return args.Get(0).(*mongo.UpdateResult), args.Error(1)
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...interface{}) *mongo.SingleResult {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(*mongo.SingleResult)
}

func TestAccountsRepository_CreateAccount(t *testing.T) {
	c := require.New(t)
	collection := new(MockCollection)
	repo := NewAccountsRepository(collection)

	account, err := domain.NewAccount("test@example.com")
	c.NoError(err)

	collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)

	err = repo.CreateAccount(context.Background(), account)

	c.NoError(err)
	collection.AssertExpectations(t)
}

func TestAccountsRepository_GetAccount(t *testing.T) {
	c := require.New(t)
	collection := new(MockCollection)
	repo := NewAccountsRepository(collection)

	account := &domain.Account{ID: "test-id"}

	t.Run("success", func(t *testing.T) {
		sr := mongo.NewSingleResultFromDocument(account, nil, bson.DefaultRegistry)
		collection.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(sr).Once()

		result, err := repo.GetAccount(context.Background(), account.ID)

		c.NoError(err)
		c.Equal(account, result)
		collection.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		sr := mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, bson.DefaultRegistry)
		collection.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(sr).Once()

		result, err := repo.GetAccount(context.Background(), account.ID)

		c.Error(err)
		c.Equal(ErrAccountNotfound, err)
		c.Nil(result)
		collection.AssertExpectations(t)
	})
}

func TestAccountsRepository_GetAccountByEmail(t *testing.T) {
	c := require.New(t)
	collection := new(MockCollection)
	repo := NewAccountsRepository(collection)

	account := &domain.Account{Email: "test@example.com"}

	t.Run("success", func(t *testing.T) {
		sr := mongo.NewSingleResultFromDocument(account, nil, bson.DefaultRegistry)
		collection.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(sr).Once()

		result, err := repo.GetAccountByEmail(context.Background(), account.Email)

		c.NoError(err)
		c.Equal(account, result)
		collection.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		sr := mongo.NewSingleResultFromDocument(nil, mongo.ErrNoDocuments, bson.DefaultRegistry)
		collection.On("FindOne", mock.Anything, mock.Anything, mock.Anything).Return(sr).Once()

		result, err := repo.GetAccountByEmail(context.Background(), account.Email)

		c.Error(err)
		c.Equal(ErrAccountNotfound, err)
		c.Nil(result)
		collection.AssertExpectations(t)
	})
}

func TestAccountsRepository_SaveGoogleAccount(t *testing.T) {
	c := require.New(t)
	collection := new(MockCollection)
	repo := NewAccountsRepository(collection)

	accountID := "test-id"
	googleAccount := &google.GoogleAccount{}

	collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)

	err := repo.SaveGoogleAccount(context.Background(), accountID, googleAccount)

	c.NoError(err)
	collection.AssertExpectations(t)
}

func TestAccountsRepository_UpdateAccount(t *testing.T) {
	c := require.New(t)
	collection := new(MockCollection)
	repo := NewAccountsRepository(collection)

	account := &domain.Account{ID: "test-id"}

	collection.On("UpdateOne", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mongo.UpdateResult{}, nil)

	err := repo.UpdateAccount(context.Background(), account)

	c.NoError(err)
	collection.AssertExpectations(t)
}
