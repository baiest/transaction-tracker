package usecase

import (
	"context"
	"transaction-tracker/internal/movements/domain"
)

// MovementUsecase define bussiness logic to manage movements.
type MovementUsecase interface {
	CreateMovement(ctx context.Context, movement *domain.Movement) error
	GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error)
	DeleteMovement(ctx context.Context, id string, accountID string) error
	GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, institutionIDs []string, limit int, offset int) (*domain.PaginatedMovements, error)
	GetMovementsByYear(ctx context.Context, accountID string, institutionIDs []string, year int) ([]*domain.Movement, error)
	GetMovementsByMonth(ctx context.Context, accountID string, institutionIDs []string, year int, month int) ([]*domain.Movement, error)
	DeleteMovementsByExtractID(ctx context.Context, extractID string) error
}
