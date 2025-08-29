package mongo

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DatabaseName string
type CollectionName string

const (
	defaultMongoURI = "mongodb://localhost:27017"

	TRANSACTIONS        DatabaseName   = "transactions"
	GOOGLE_ACCOUNTS     CollectionName = "google_accounts"
	GMAIL_EXTRACTS      CollectionName = "gmail_extracts"
	GMAIL_NOTIFICATIONS CollectionName = "gmail_notifications"
	GMAIL_MESSAGES      CollectionName = "gmail_messages"
)

type MongoClient struct {
	client *mongo.Client
}

// NewClient returns a connected *mongo.Client using MONGO_URI if present.
func NewClient(ctx context.Context) (*MongoClient, error) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = defaultMongoURI
	}

	_client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("error creating mongo client: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return &MongoClient{client: _client}, nil
}

func (c *MongoClient) Collection(db DatabaseName, collection CollectionName) (*mongo.Collection, error) {
	if c.client == nil {
		return nil, fmt.Errorf("client not initialized")
	}

	return c.client.Database(string(db)).Collection(string(collection)), nil
}
