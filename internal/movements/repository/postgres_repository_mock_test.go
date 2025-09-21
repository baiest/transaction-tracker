package repository

import (
	"context"
	"testing"
	"transaction-tracker/internal/movements/domain"

	"github.com/stretchr/testify/assert"
)

func TestMockMovementRepository_CreateMovement(t *testing.T) {
	mockRepo := new(MockMovementRepository)
	ctx := context.Background()

	movement := &domain.Movement{ID: "mov1"}

	mockRepo.On("CreateMovement", ctx, movement).Return(nil)

	err := mockRepo.CreateMovement(ctx, movement)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestMockMovementRepository_GetMovementByID(t *testing.T) {
	mockRepo := new(MockMovementRepository)
	ctx := context.Background()

	expectedMovement := &domain.Movement{ID: "mov1"}

	mockRepo.On("GetMovementByID", ctx, "mov1").Return(expectedMovement, nil)

	result, err := mockRepo.GetMovementByID(ctx, "mov1")

	assert.NoError(t, err)
	assert.Equal(t, expectedMovement, result)
	mockRepo.AssertExpectations(t)
}

func TestMockMovementRepository_GetMovementsByAccountID(t *testing.T) {
	mockRepo := new(MockMovementRepository)
	ctx := context.Background()

	expectedMovements := []*domain.Movement{
		{ID: "mov1"},
		{ID: "mov2"},
	}

	mockRepo.On("GetMovementsByAccountID", ctx, "acc1", 10, 0).Return(expectedMovements, nil)

	result, err := mockRepo.GetMovementsByAccountID(ctx, "acc1", 10, 0)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedMovements, result)
	mockRepo.AssertExpectations(t)
}

func TestMockMovementRepository_GetTotalMovementsByAccountID(t *testing.T) {
	mockRepo := new(MockMovementRepository)
	ctx := context.Background()

	mockRepo.On("GetTotalMovementsByAccountID", ctx, "acc1").Return(5, nil)

	total, err := mockRepo.GetTotalMovementsByAccountID(ctx, "acc1")

	assert.NoError(t, err)
	assert.Equal(t, 5, total)
	mockRepo.AssertExpectations(t)
}
