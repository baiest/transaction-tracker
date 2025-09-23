package repository

import (
	"context"

	"transaction-tracker/internal/movements/domain"

	"github.com/stretchr/testify/mock"
)

// MockMovementRepository is a mock of the repository interface.
type MockMovementRepository struct {
	mock.Mock
}

// CreateMovement simulates the creation of a movement.
func (m *MockMovementRepository) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	args := m.Called(ctx, movement)
	return args.Error(0)
}

// GetMovementByID simulates retrieving a movement by its ID.
func (m *MockMovementRepository) GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error) {
	args := m.Called(ctx, id, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Movement), args.Error(1)
}

// GetMovementsByAccountID simulates retrieving movements by account ID.
func (m *MockMovementRepository) GetMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Movement, error) {
	args := m.Called(ctx, accountID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]*domain.Movement), args.Error(1)
}

// GetTotalMovementsByAccountID simulates retrieving the total number of movements for an account ID.
func (m *MockMovementRepository) GetTotalMovementsByAccountID(ctx context.Context, accountID string) (int, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return 0, args.Error(1)
	}

	return args.Get(0).(int), args.Error(1)
}
