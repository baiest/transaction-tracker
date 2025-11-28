package repository

import (
	"context"
	"transaction-tracker/internal/movements/domain"
)

// MovementRepository defines the contract for a data store to handle movements.
type MovementRepository interface {
	CreateMovement(ctx context.Context, movement *domain.Movement) error
	GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error)
	Delete(ctx context.Context, id string, accountID string) error
	GetTotalMovementsByAccountID(ctx context.Context, accountID string) (int, error)
	GetMovementsByAccountID(ctx context.Context, accountID string, institutionIDs []string, limit int, offset int) ([]*domain.Movement, error)
	DeleteMovementsByExtractID(ctx context.Context, extractID string) error
}
