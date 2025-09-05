package repositories

import (
	"context"
	"time"
	databaseMongo "transaction-tracker/database/mongo"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	movementsDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	movementsCollection databaseMongo.CollectionName = databaseMongo.MOVEMENTS
)

type IMovementsRepository interface {
	SaveMovement(context.Context, *schemas.Movement) error
	GetMovements(context.Context, int64) ([]*schemas.Movement, int64, error)
	GetMovementsByDateRange(context.Context, time.Time, time.Time) ([]*schemas.Movement, error)
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

func (r *MovementsRepository) GetMovementsByDateRange(ctx context.Context, startDate time.Time, finishDate time.Time) ([]*schemas.Movement, error) {
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

func (r *MovementsRepository) GetMovements(ctx context.Context, page int64) ([]*schemas.Movement, int64, error) {
	var movements []*schemas.Movement

	repo, err := NewMovementsRepository(ctx)
	if err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}

	skip := (page - 1) * 10

	filter := bson.M{}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	totalPages := (total + 10 - 1) / 10

	findOptions := options.Find().
		SetLimit(10).
		SetSkip(skip).
		SetSort(bson.M{"date": -1})

	cursor, err := repo.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var movement schemas.Movement

		if err := cursor.Decode(&movement); err != nil {
			return nil, 0, err
		}

		movements = append(movements, &movement)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return movements, totalPages, nil
}
