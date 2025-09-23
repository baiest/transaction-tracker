package usecase

import (
	"context"
	"transaction-tracker/internal/movements/domain"
)

// MovementUsecase define bussiness logic to manage movements.
type MovementUsecase interface {
	CreateMovement(ctx context.Context, movement *domain.Movement) error
	GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error)
	GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) (*domain.PaginatedMovements, error)
}
