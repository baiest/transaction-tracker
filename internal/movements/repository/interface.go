package repository

import (
	"context"
	"transaction-tracker/internal/movements/domain"
)

// MovementRepository defines the contract for a data store to handle movements.
type MovementRepository interface {
	CreateMovement(ctx context.Context, movement *domain.Movement) error
	GetMovementByID(ctx context.Context, id string) (*domain.Movement, error)
	GetTotalMovementsByAccountID(ctx context.Context, accountID string) (int, error)
	GetMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Movement, error)
}
