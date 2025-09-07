package repositories

import (
	"context"
	"errors"
	"time"
	databaseMongo "transaction-tracker/database/mongo"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	accountsDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	accountsCollection databaseMongo.CollectionName = databaseMongo.ACCOUNTS
)

type IAccountsRepository interface {
	CreateAccount(context.Context, *schemas.Account) error
	GetAccount(context.Context, string) (*schemas.Account, error)
	GetAccountByEmail(context.Context, string) (*schemas.Account, error)
}

type AccountsRepository struct {
	collection *mongo.Collection
}

var (
	ErrAccountNotfound = errors.New("account not found")
)

func NewAccountsRepository(ctx context.Context) (*AccountsRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(accountsDatabase, accountsCollection)
	if err != nil {
		return nil, err
	}

	return &AccountsRepository{collection: collection}, nil
}

func (a *AccountsRepository) CreateAccount(ctx context.Context, account *schemas.Account) error {
	_, err := a.collection.UpdateOne(
		ctx,
		bson.M{"_id": account.ID},
		bson.M{
			"$set": bson.M{
				"email":         account.Email,
				"refresh_token": account.RefreshToken,
				"created_at":    time.Now(),
				"updated_at":    time.Now(),
			},
		},
		options.UpdateOne().SetUpsert(true),
	)

	return err
}

func (r *AccountsRepository) GetAccount(ctx context.Context, accountID string) (*schemas.Account, error) {
	result := r.collection.FindOne(ctx, bson.M{"_id": accountID})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrAccountNotfound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	account := &schemas.Account{}
	if err := result.Decode(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (r *AccountsRepository) GetAccountByEmail(ctx context.Context, email string) (*schemas.Account, error) {
	var account *schemas.Account

	result := r.collection.FindOne(ctx, bson.M{"email": email})

	if err := result.Decode(&account); err != nil {
		return nil, err
	}

	return account, nil
}
