package repository

import (
	"errors"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	databaseMongo "transaction-tracker/database/mongo"
)

var (
	ErrNotificationAlreadyExists = errors.New("gmail notification already exists")
)

const (
	gmailNotificationsDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	gmailNotificationsCollection databaseMongo.CollectionName = databaseMongo.GMAIL_NOTIFICATIONS
)

type IGmailNotificationsRepository interface {
	SaveNotification(ctx Context, notification *schemas.GmailNotification) error
}

type GmailNotificationsRepository struct {
	collection *mongo.Collection
}

func NewGmailNotificationsRepository(ctx Context) (*GmailNotificationsRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(gmailNotificationsDatabase, gmailNotificationsCollection)
	if err != nil {
		return nil, err
	}

	return &GmailNotificationsRepository{collection: collection}, nil
}

// SaveNotification inserts a new GmailNotification. If the ID already exists, returns ErrNotificationAlreadyExists.
func (r *GmailNotificationsRepository) SaveNotification(ctx Context, notification *schemas.GmailNotification) error {
	_, err := r.collection.InsertOne(ctx, notification, options.InsertOne())
	if err != nil {
		// Check for duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return ErrNotificationAlreadyExists
		}

		return err
	}

	return nil
}

// GetNotificationByID fetches a GmailNotification by its ID.
func (r *GmailNotificationsRepository) GetNotificationByID(ctx Context, id string) (*schemas.GmailNotification, error) {
	result := r.collection.FindOne(ctx, bson.M{"_id": id})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, errors.New("gmail notification not found")
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	notification := &schemas.GmailNotification{}
	if err := result.Decode(notification); err != nil {
		return nil, err
	}

	return notification, nil
}

// UpdateNotification updates an existing GmailNotification by its ID.
func (r *GmailNotificationsRepository) UpdateNotification(ctx Context, notification *schemas.GmailNotification) error {
	filter := bson.M{"_id": notification.ID}
	update := bson.M{
		"$set": bson.M{
			"email":  notification.Email,
			"status": notification.Status,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("gmail notification not found")
	}

	return nil
}
