package usecase

import (
	"context"
	"transaction-tracker/internal/movements/domain"
)

// MockMovementUsecase is a mock implementation of the MovementUsecase interface.
type MockMovementUsecase struct {
	CreateMovementFunc                   func(ctx context.Context, movement *domain.Movement) error
	GetPaginatedMovementsByAccountIDFunc func(ctx context.Context, accountID string, limit int, offset int) (*domain.PaginatedMovements, error)
	GetMovementByIDFunc                  func(ctx context.Context, id string) (*domain.Movement, error)
}

// CreateMovement calls the mocked CreateMovementFunc.
func (m *MockMovementUsecase) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	return m.CreateMovementFunc(ctx, movement)
}

// GetPaginatedMovementsByAccountID calls the mocked GetPaginatedMovementsByAccountIDFunc.
func (m *MockMovementUsecase) GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) (*domain.PaginatedMovements, error) {
	return m.GetPaginatedMovementsByAccountIDFunc(ctx, accountID, limit, offset)
}

// GetMovementByID calls the mocked GetMovementByIDFunc.
func (m *MockMovementUsecase) GetMovementByID(ctx context.Context, id string) (*domain.Movement, error) {
	return m.GetMovementByIDFunc(ctx, id)
}
