package repositories

import (
	"errors"
	"transaction-tracker/database/mongo/schemas"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	databaseMongo "transaction-tracker/database/mongo"
)

var (
	ErrMessageAlreadyExists = errors.New("gmail message already exists")
	ErrMessageNotFound      = errors.New("gmail message not found")
)

const (
	gmailMessagesDatabase   databaseMongo.DatabaseName   = databaseMongo.TRANSACTIONS
	gmailMessagesCollection databaseMongo.CollectionName = databaseMongo.GMAIL_MESSAGES
)

type IGmailMessageRepository interface {
	SaveMessage(ctx Context, message *schemas.Message) error
	GetMessageByID(ctx Context, id string) (*schemas.Message, error)
	GetMessagesByNotificationID(ctx Context, notificationID string) ([]*schemas.Message, error)
	UpdateMessage(ctx Context, message *schemas.Message) error
}

type GmailMessageRepository struct {
	collection *mongo.Collection
}

func NewGmailMessageRepository(ctx Context) (*GmailMessageRepository, error) {
	client, err := databaseMongo.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	collection, err := client.Collection(gmailMessagesDatabase, gmailMessagesCollection)
	if err != nil {
		return nil, err
	}

	return &GmailMessageRepository{collection: collection}, nil
}

// SaveMessage inserts a new Message. If the ID already exists, returns ErrMessageAlreadyExists.
func (r *GmailMessageRepository) SaveMessage(ctx Context, message *schemas.Message) error {
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

// GetMessageByID fetches a Message by its ID.
func (r *GmailMessageRepository) GetMessageByID(ctx Context, id string) (*schemas.Message, error) {
	result := r.collection.FindOne(ctx, bson.M{"_id": id})
	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrMessageNotFound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	message := &schemas.Message{}
	if err := result.Decode(message); err != nil {
		return nil, err
	}

	return message, nil
}

// GetMessagesByNotificationID fetches all Messages by their notification ID.
func (r *GmailMessageRepository) GetMessagesByNotificationID(ctx Context, notificationID string) ([]*schemas.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"notification_id": notificationID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var messages []*schemas.Message
	if err = cursor.All(ctx, &messages); err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return []*schemas.Message{}, nil
	}

	return messages, nil
}

// UpdateMessage updates an existing Message by its ID.
func (r *GmailMessageRepository) UpdateMessage(ctx Context, message *schemas.Message) error {
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
