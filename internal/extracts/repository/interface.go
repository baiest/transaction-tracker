package repository

import (
	"context"

	"transaction-tracker/internal/extracts/domain"
)

// ExtractsRepository defines the interface for extract data persistence.
type ExtractsRepository interface {
	GetByMessageID(ctx context.Context, messageID string) (*domain.Extract, error)
	Save(ctx context.Context, extract *domain.Extract) error
	Update(ctx context.Context, extract *domain.Extract) error
}
