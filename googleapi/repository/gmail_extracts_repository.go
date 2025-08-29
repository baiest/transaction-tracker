package repository

import (
	"errors"
	"transaction-tracker/database/mongo/schemas"

	databaseMongo "transaction-tracker/database/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	gmailExtractsDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	gmailExtractsCollection databaseMongo.CollectionName = databaseMongo.GMAIL_EXTRACTS
)

var (
	ErrExtractAlreadyExists = errors.New("extract already exists")
	ErrExtractNotFound      = errors.New("extract not found")
)

type IGmailExtractsRepository interface {
	SaveExtract(ctx Context, extract *schemas.GmailExtract) error
	GetExtract(ctx Context, id string) (*schemas.GmailExtract, error)
}

type GmailExtractsRepository struct {
	collection *mongo.Collection
}

func NewGmailExtractsRepository(ctx Context) (*GmailExtractsRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(gmailExtractsDatabase, gmailExtractsCollection)
	if err != nil {
		return nil, err
	}

	return &GmailExtractsRepository{collection: collection}, nil
}

// SaveExtract inserts a new Extract. If the ID already exists, returns ErrExtractAlreadyExists.
func (r *GmailExtractsRepository) SaveExtract(ctx Context, extract *schemas.GmailExtract) error {
	_, err := r.collection.InsertOne(ctx, extract, options.InsertOne())
	if err != nil {
		// Check for duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return ErrExtractAlreadyExists
		}

		return err
	}

	return nil
}

func (r *GmailExtractsRepository) GetExtract(ctx Context, id string) (*schemas.GmailExtract, error) {
	var extract schemas.GmailExtract

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&extract)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, ErrExtractNotFound
	}

	if err != nil {
		return nil, err
	}

	return &extract, nil
}
