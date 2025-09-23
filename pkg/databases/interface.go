package databases

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Client is an interface for our database client
// that the rest of the application will use.
type Client interface {
	GetPool() *pgxpool.Pool
	Close()
}

// CollectionAPI abstracts methods we use from mongo.Collection
type CollectionAPI interface {
	Find(ctx context.Context, filter any, opts ...options.Lister[options.FindOptions]) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter any, opts ...options.Lister[options.FindOneOptions]) *mongo.SingleResult
	InsertOne(ctx context.Context, doc any, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter any, update any, opts ...options.Lister[options.UpdateOneOptions]) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter any, opts ...options.Lister[options.DeleteOneOptions]) (*mongo.DeleteResult, error)
}
