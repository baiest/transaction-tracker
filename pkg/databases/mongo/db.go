package mongo

import (
	"context"
	"fmt"
	"os"
	"time"
	"transaction-tracker/pkg/databases"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DatabaseName string
type CollectionName string

const (
	defaultMongoURI = "mongodb://localhost:27017"

	TRANSACTIONS DatabaseName   = "transactions"
	EXTRACTS     CollectionName = "extracts"
	MESSAGES     CollectionName = "messages"
	MOVEMENTS    CollectionName = "movements"
	ACCOUNTS     CollectionName = "accounts"
)

type MongoClient struct {
	client *mongo.Client
}

// NewClient returns a connected *mongo.Client using MONGO_URI if present.
func NewClient(ctx context.Context) (context.Context, *MongoClient, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = defaultMongoURI
	}

	opts := options.Client().
		ApplyURI(uri).
		SetMinPoolSize(5).
		SetMaxPoolSize(50).
		SetMaxConnIdleTime(30 * time.Second)

	_client, err := mongo.Connect(opts)
	if err != nil {
		return ctx, nil, fmt.Errorf("error creating mongo client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := _client.Ping(pingCtx, nil); err != nil {
		return ctx, nil, fmt.Errorf("mongo ping failed: %w", err)
	}

	return ctx, &MongoClient{client: _client}, nil
}

func (c *MongoClient) Collection(db DatabaseName, collection CollectionName) (databases.CollectionAPI, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return c.client.Database(string(db)).Collection(string(collection)), nil
}
