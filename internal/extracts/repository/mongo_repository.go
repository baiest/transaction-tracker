package repository

import (
	"context"
	"errors"
	"time"
	"transaction-tracker/internal/extracts/domain"
	"transaction-tracker/pkg/databases"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrExtractNotFound = errors.New("extract not found")
)

type extractsRepository struct {
	collection databases.CollectionAPI
	nowFunc    func() time.Time
}

// NewExtractsRepository creates a new instance of ExtractsRepository.
func NewExtractsRepository(collection databases.CollectionAPI) ExtractsRepository {
	return &extractsRepository{
		collection: collection,
		nowFunc:    time.Now,
	}
}

func (r *extractsRepository) GetByMessageID(ctx context.Context, messageID string) (*domain.Extract, error) {
	result := r.collection.FindOne(ctx, map[string]interface{}{
		"message_id": messageID,
	})

	if errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return nil, ErrExtractNotFound
	}

	if result.Err() != nil {
		return nil, result.Err()
	}

	extract := &domain.Extract{}
	if err := result.Decode(extract); err != nil {
		return nil, err
	}

	return extract, nil
}

// Save saves the given extract to the database.
func (r *extractsRepository) Save(ctx context.Context, extract *domain.Extract) error {
	now := r.nowFunc()
	extract.CreatedAt = now
	extract.UpdatedAt = now

	_, err := r.collection.InsertOne(ctx, extract)

	return err
}

// Update updates the given extract in the database.
func (r *extractsRepository) Update(ctx context.Context, extract *domain.Extract) error {
	extract.UpdatedAt = r.nowFunc()

	filter := map[string]interface{}{
		"_id": extract.ID,
	}

	update := map[string]interface{}{
		"$set": extract,
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}
