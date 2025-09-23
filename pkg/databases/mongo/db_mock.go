package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockCollection struct {
	InsertOneFn func(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	FindOneFn   func(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult
	FindFn      func(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error)
	UpdateOneFn func(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)
}

// MockMongoClient is a fake implementation of MongoClient for testing purposes.
type MockMongoClient struct {
	Collections map[string]*mongo.Collection
	Fail        bool
}

func NewMockMongoClient() *MockMongoClient {
	return &MockMongoClient{
		Collections: make(map[string]*mongo.Collection),
	}
}

func (m *MockMongoClient) Collection(db DatabaseName, collection CollectionName) (*mongo.Collection, error) {
	if m.Fail {
		return nil, fmt.Errorf("forced error")
	}

	key := fmt.Sprintf("%s.%s", db, collection)
	if col, ok := m.Collections[key]; ok {
		return col, nil
	}

	// Return a dummy mongo.Collection pointer (not connected to a server).
	dummy := &mongo.Collection{}
	m.Collections[key] = dummy
	return dummy, nil
}

func (m *MockCollection) InsertOne(ctx context.Context, doc any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.InsertOneFn(ctx, doc, opts...)
}

func (m *MockCollection) FindOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.FindOneFn(ctx, filter, opts...)
}

func (m *MockCollection) Find(ctx context.Context, filter any, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.FindFn(ctx, filter, opts...)
}

func (m *MockCollection) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return m.UpdateOneFn(ctx, filter, update, opts...)
}
