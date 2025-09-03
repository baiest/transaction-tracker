package repositories

import (
	"context"
	"time"
	databaseMongo "transaction-tracker/database/mongo"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	movementsDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	movementsCollection databaseMongo.CollectionName = databaseMongo.MOVEMENTS
)

type IMovementsRepository interface {
	SaveMovement(context.Context, *schemas.Movement) error
	GetMovements(context.Context, time.Time, time.Time) ([]*schemas.Movement, error)
}

type MovementsRepository struct {
	collection *mongo.Collection
}

func NewMovementsRepository(ctx context.Context) (*MovementsRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(movementsDatabase, movementsCollection)
	if err != nil {
		return nil, err
	}

	return &MovementsRepository{collection: collection}, nil
}

func (r *MovementsRepository) SaveMovement(ctx context.Context, movement *schemas.Movement) error {
	_, err := r.collection.InsertOne(ctx, movement)

	return err
}

func (r *MovementsRepository) GetMovements(ctx context.Context, startDate time.Time, finishDate time.Time) ([]*schemas.Movement, error) {
	var movements []*schemas.Movement

	cursor, err := r.collection.Find(ctx, bson.M{"date": bson.M{"$gte": startDate, "$lte": finishDate}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var movement *schemas.Movement

		if err := cursor.Decode(&movement); err != nil {
			return nil, err
		}

		movements = append(movements, movement)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}
