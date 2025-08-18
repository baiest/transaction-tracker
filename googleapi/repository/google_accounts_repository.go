package repository

import (
	"errors"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	databaseMongo "transaction-tracker/database/mongo"
)

// IGoogleAccountsRepository defines the operations to create and query Google tokens.
type IGoogleAccountsRepository interface {
	// SaveToken creates or updates a token by EmailAddress (upsert semantics).
	SaveToken(ctx Context, token *schemas.GoogleAccount) error

	// GetTokenByEmail fetches the token for the given EmailAddress.
	GetTokenByEmail(ctx Context, emailAddress string) (*schemas.GoogleAccount, error)
}

const (
	database   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	collection databaseMongo.CollectionName = databaseMongo.GOOGLE_ACCOUNTS
)

var (
	ErrTokenNotFound = errors.New("token not found")
)

// GoogleAccountsRepository implements TokenRepository using a MongoDB collection.
type GoogleAccountsRepository struct {
	collection *mongo.Collection
}

// NeGoogleAccountsRepository creates a new GoogleAccountsRepository with the provided collection.
func NeGoogleAccountsRepository(ctx Context) (*GoogleAccountsRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(database, collection)
	if err != nil {
		return nil, err
	}

	return &GoogleAccountsRepository{collection: collection}, nil
}

// SaveToken upserts the token by EmailAddress.
func (r *GoogleAccountsRepository) SaveToken(ctx Context, account *schemas.GoogleAccount) error {

	_, err := r.collection.InsertOne(ctx, account)
	return err
}

// GetTokenByEmail fetches a token by EmailAddress.
func (r *GoogleAccountsRepository) GetTokenByEmail(ctx Context, emailAddress string) (*schemas.GoogleAccount, error) {
	result := r.collection.FindOne(ctx, bson.M{"_id": emailAddress})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrTokenNotFound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	account := &schemas.GoogleAccount{}
	if err := result.Decode(account); err != nil {
		return nil, err
	}

	return account, nil
}
