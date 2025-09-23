package usecase

import (
	"context"
	"errors"
	"testing"
	"transaction-tracker/internal/movements/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockMovementUsecase_CreateMovement(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	movement := &domain.Movement{ID: "mov1"}

	t.Run("success", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			CreateMovementFunc: func(cx context.Context, m *domain.Movement) error {
				c.Equal(ctx, cx)
				c.Equal(movement, m)
				return nil
			},
		}

		err := mockUC.CreateMovement(ctx, movement)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			CreateMovementFunc: func(c context.Context, m *domain.Movement) error {
				return errors.New("create failed")
			},
		}

		err := mockUC.CreateMovement(ctx, movement)
		assert.EqualError(t, err, "create failed")
	})
}

func TestMockMovementUsecase_GetMovementsByAccountID(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	accountID := "acc1"
	movements := []*domain.Movement{{ID: "mov1"}}

	t.Run("success", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			GetPaginatedMovementsByAccountIDFunc: func(cx context.Context, acc string, limit int, offset int) (*domain.PaginatedMovements, error) {
				c.Equal(ctx, cx)
				c.Equal(accountID, acc)
				c.Equal(accountID, acc)
				c.Equal(accountID, acc)

				return &domain.PaginatedMovements{
					Movements: movements,
				}, nil
			},
		}

		result, err := mockUC.GetPaginatedMovementsByAccountIDFunc(ctx, accountID, 10, 1)
		assert.NoError(t, err)

		c.Equal(movements, result.Movements)
	})

	t.Run("error", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			GetPaginatedMovementsByAccountIDFunc: func(ctx context.Context, accountID string, limit, offset int) (*domain.PaginatedMovements, error) {
				return nil, errors.New("fetch failed")
			},
		}

		result, err := mockUC.GetPaginatedMovementsByAccountIDFunc(ctx, accountID, 10, 1)
		assert.Nil(t, result)
		assert.EqualError(t, err, "fetch failed")
	})
}

func TestMockMovementUsecase_GetMovementByID(t *testing.T) {
	c := require.New(t)

	ctx := context.Background()
	movement := &domain.Movement{ID: "mov1"}

	t.Run("success", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			GetMovementByIDFunc: func(cx context.Context, id string, accountID string) (*domain.Movement, error) {
				c.Equal(ctx, cx)
				c.Equal("mov1", id)
				c.Equal("acc1", accountID)
				return movement, nil
			},
		}

		result, err := mockUC.GetMovementByID(ctx, "mov1", "acc1")
		assert.NoError(t, err)
		c.Equal(movement, result)
	})

	t.Run("error", func(t *testing.T) {
		mockUC := &MockMovementUsecase{
			GetMovementByIDFunc: func(c context.Context, id string, accountID string) (*domain.Movement, error) {
				return nil, errors.New("not found")
			},
		}

		result, err := mockUC.GetMovementByID(ctx, "mov1", "acc1")
		assert.Nil(t, result)
		assert.EqualError(t, err, "not found")
	})
}
