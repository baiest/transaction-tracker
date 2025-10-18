package usecase

import (
	"context"
	"transaction-tracker/internal/movements/domain"

	"github.com/stretchr/testify/mock"
)

type MockMovementUsecase struct {
	mock.Mock
}

func (m *MockMovementUsecase) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	if m == nil {
		return nil
	}
	args := m.Called(ctx, movement)
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockMovementUsecase) GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error) {
	if m == nil {
		return nil, nil
	}
	args := m.Called(ctx, id, accountID)

	var movement *domain.Movement
	if val := args.Get(0); val != nil {
		if cast, ok := val.(*domain.Movement); ok {
			movement = cast
		}
	}
	return movement, args.Error(1)
}

func (m *MockMovementUsecase) GetMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Movement, error) {
	if m == nil {
		return nil, nil
	}
	args := m.Called(ctx, accountID, limit, offset)

	var movements []*domain.Movement
	if val := args.Get(0); val != nil {
		if cast, ok := val.([]*domain.Movement); ok {
			movements = cast
		}
	}
	return movements, args.Error(1)
}

func (m *MockMovementUsecase) GetTotalMovementsByAccountID(ctx context.Context, accountID string) (int, error) {
	if m == nil {
		return 0, nil
	}
	args := m.Called(ctx, accountID)
	return args.Int(0), args.Error(1)
}

func (m *MockMovementUsecase) Delete(ctx context.Context, id string, accountID string) error {
	if m == nil {
		return nil
	}
	args := m.Called(ctx, id, accountID)
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockMovementUsecase) DeleteMovementsByExtractID(ctx context.Context, extractID string) error {
	if m == nil {
		return nil
	}
	args := m.Called(ctx, extractID)
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockMovementUsecase) DeleteMovement(ctx context.Context, id string, accountID string) error {
	if m == nil {
		return nil
	}
	args := m.Called(ctx, id, accountID)
	if len(args) == 0 {
		return nil
	}
	return args.Error(0)
}

func (m *MockMovementUsecase) GetMovementsByYear(ctx context.Context, accountID string, year int) ([]*domain.Movement, error) {
	if m == nil {
		return nil, nil
	}
	args := m.Called(ctx, accountID, year)

	var movements []*domain.Movement
	if val := args.Get(0); val != nil {
		if cast, ok := val.([]*domain.Movement); ok {
			movements = cast
		}
	}
	return movements, args.Error(1)
}

func (m *MockMovementUsecase) GetMovementsByMonth(ctx context.Context, accountID string, year int, month int) ([]*domain.Movement, error) {
	if m == nil {
		return nil, nil
	}
	args := m.Called(ctx, accountID, year, month)

	var movements []*domain.Movement
	if val := args.Get(0); val != nil {
		if cast, ok := val.([]*domain.Movement); ok {
			movements = cast
		}
	}
	return movements, args.Error(1)
}

func (m *MockMovementUsecase) GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) (*domain.PaginatedMovements, error) {
	if m == nil {
		return nil, nil
	}

	args := m.Called(ctx, accountID, limit, offset)

	var movements *domain.PaginatedMovements
	if val := args.Get(0); val != nil {
		if cast, ok := val.(*domain.PaginatedMovements); ok {
			movements = cast
		}
	}

	return movements, args.Error(1)
}
