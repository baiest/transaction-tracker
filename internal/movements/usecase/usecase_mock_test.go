package usecase

import (
	"context"
	"errors"
	"testing"
	"transaction-tracker/internal/movements/domain"

	"github.com/stretchr/testify/require"
)

func TestMockCreateMovement(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	movement := &domain.Movement{ID: "123"}

	mockRepo.On("CreateMovement", ctx, movement).Return(nil)

	err := mockRepo.CreateMovement(ctx, movement)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateMovement_WithError(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	movement := &domain.Movement{ID: "999"}

	mockRepo.On("CreateMovement", ctx, movement).Return(errors.New("db error"))

	err := mockRepo.CreateMovement(ctx, movement)
	require.EqualError(t, err, "db error")

	mockRepo.AssertExpectations(t)
}

func TestCreateMovement_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase // intentionally nil
	err := mockRepo.CreateMovement(context.Background(), nil)
	require.NoError(t, err)
}

func TestGetMovementByID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	expected := &domain.Movement{ID: "123"}

	mockRepo.On("GetMovementByID", ctx, "123", "acc-001").Return(expected, nil)

	result, err := mockRepo.GetMovementByID(ctx, "123", "acc-001")
	require.NoError(t, err)
	require.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}

func TestGetMovementByID_NilReturn(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	mockRepo.On("GetMovementByID", ctx, "notfound", "acc-001").Return(nil, errors.New("not found"))

	result, err := mockRepo.GetMovementByID(ctx, "notfound", "acc-001")
	require.Nil(t, result)
	require.EqualError(t, err, "not found")

	mockRepo.AssertExpectations(t)
}

func TestGetMovementByID_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase
	result, err := mockRepo.GetMovementByID(context.Background(), "id", "acc")
	require.Nil(t, result)
	require.NoError(t, err)
}

func TestGetMovementsByAccountID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	movements := []*domain.Movement{
		{ID: "1"},
		{ID: "2"},
	}

	mockRepo.On("GetMovementsByAccountID", ctx, "acc-001", []string{}, 10, 0).Return(movements, nil)

	result, err := mockRepo.GetMovementsByAccountID(ctx, "acc-001", []string{}, 10, 0)
	require.NoError(t, err)
	require.Equal(t, movements, result)

	mockRepo.AssertExpectations(t)
}

func TestGetMovementsByAccountID_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase
	result, err := mockRepo.GetMovementsByAccountID(context.Background(), "acc", []string{}, 10, 0)
	require.Nil(t, result)
	require.NoError(t, err)
}

func TestGetTotalMovementsByAccountID(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	mockRepo.On("GetTotalMovementsByAccountID", ctx, "acc-001").Return(5, nil)

	count, err := mockRepo.GetTotalMovementsByAccountID(ctx, "acc-001")
	require.NoError(t, err)
	require.Equal(t, 5, count)

	mockRepo.AssertExpectations(t)
}

func TestGetTotalMovementsByAccountID_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase
	count, err := mockRepo.GetTotalMovementsByAccountID(context.Background(), "acc")
	require.Equal(t, 0, count)
	require.NoError(t, err)
}

func TestDelete_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	mockRepo.On("Delete", ctx, "1", "acc-001").Return(nil)

	err := mockRepo.Delete(ctx, "1", "acc-001")
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDelete_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase
	err := mockRepo.Delete(context.Background(), "1", "acc")
	require.NoError(t, err)
}

func TestDeleteMovementsByExtractID_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockMovementUsecase)
	mockRepo.On("DeleteMovementsByExtractID", ctx, "extract-001").Return(nil)

	err := mockRepo.DeleteMovementsByExtractID(ctx, "extract-001")
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteMovementsByExtractID_NilRepo(t *testing.T) {
	var mockRepo *MockMovementUsecase
	err := mockRepo.DeleteMovementsByExtractID(context.Background(), "extract")
	require.NoError(t, err)
}
