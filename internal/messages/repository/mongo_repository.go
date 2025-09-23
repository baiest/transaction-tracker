package repository

import (
	"context"
	"errors"
	"time"
	"transaction-tracker/internal/messages/domain"

	"transaction-tracker/pkg/databases"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var (
	ErrMessageAlreadyExists = errors.New("gmail message already exists")
	ErrMessageNotFound      = errors.New("gmail message not found")
)

type messageRepository struct {
	collection databases.CollectionAPI
	nowFunc    func() time.Time
}

func NewMessageRepository(ctx context.Context, collection databases.CollectionAPI) MessageRepository {
	return &messageRepository{collection: collection, nowFunc: time.Now}
}

// SaveMessage inserts a new Message. If the ID already exists, returns ErrMessageAlreadyExists.
func (r *messageRepository) SaveMessage(ctx context.Context, message *domain.Message) error {
	now := r.nowFunc()

	message.CreatedAt = now
	message.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, message, options.InsertOne())
	if err != nil {
		// Check for duplicate key error
		if mongo.IsDuplicateKeyError(err) {
			return ErrMessageAlreadyExists
		}

		return err
	}

	return nil
}

func (r *messageRepository) GetMessageByExternalID(ctx context.Context, id string, accountID string) (*domain.Message, error) {
	filter := bson.M{
		"external_id": id,
		"account_id":  accountID,
	}

	result := r.collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrMessageNotFound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	message := &domain.Message{}
	if err := result.Decode(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessageByID fetches a Message by its ID and AccountID.
func (r *messageRepository) GetMessageByID(ctx context.Context, id, accountID string) (*domain.Message, error) {
	filter := bson.M{
		"_id":        id,
		"account_id": accountID,
	}

	result := r.collection.FindOne(ctx, filter)
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrMessageNotFound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	message := &domain.Message{}
	if err := result.Decode(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessagesByNotificationID fetches all Messages by their notification ID.
func (r *messageRepository) GetMessagesByNotificationID(ctx context.Context, notificationID string) ([]*domain.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"notification_id": notificationID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var messages []*domain.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return []*domain.Message{}, nil
	}

	return messages, nil
}

// UpdateMessage updates an existing Message by its ID.
func (r *messageRepository) UpdateMessage(ctx context.Context, message *domain.Message) error {
	filter := bson.M{"_id": message.ID}
	update := bson.M{
		"$set": bson.M{
			"status": message.Status,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("gmail message not found")
	}

	return nil
}
