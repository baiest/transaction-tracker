package repository

import (
	"context"
	"errors"
	"time"
	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/pkg/databases"
	"transaction-tracker/pkg/google"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type accountsRepository struct {
	collection databases.CollectionAPI
	nowFunc    func() time.Time
}

var (
	ErrAccountNotfound = errors.New("account not found")
)

// NewAccountsRepository creates a new AccountsRepository with the provided collection.
func NewAccountsRepository(collection databases.CollectionAPI) AccountsRepository {
	return &accountsRepository{collection: collection, nowFunc: time.Now}
}

func (a *accountsRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	_, err := a.collection.UpdateOne(
		ctx,
		bson.M{"_id": account.ID},
		bson.M{
			"$set": bson.M{
				"email":         account.Email,
				"refresh_token": account.RefreshToken,
				"created_at":    a.nowFunc(),
				"updated_at":    a.nowFunc(),
			},
		},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (a *accountsRepository) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	result := a.collection.FindOne(ctx, bson.M{"_id": accountID})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrAccountNotfound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	account := &domain.Account{}
	if err := result.Decode(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (a *accountsRepository) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	var account *domain.Account

	result := a.collection.FindOne(ctx, bson.M{"email": email})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrAccountNotfound
	}

	if err := result.Decode(&account); err != nil {
		return nil, err
	}

	return account, nil
}

// SaveToken upserts the token by EmailAddress.
func (a *accountsRepository) SaveGoogleAccount(ctx context.Context, accountID string, googleAccount *google.GoogleAccount) error {
	_, err := a.collection.UpdateOne(
		ctx,
		bson.M{"_id": accountID},
		bson.M{
			"$set": bson.M{
				"google_account": googleAccount,
				"updated_at":     time.Now(),
			},
		},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (a *accountsRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	_, err := a.collection.UpdateOne(
		ctx,
		bson.M{"_id": account.ID},
		bson.M{
			"$set": bson.M{
				"email":      account.Email,
				"updated_at": a.nowFunc(),
			},
		},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}
