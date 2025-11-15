package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/repository"
)

func TestCreateMovement(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)

	u := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		movement := &domain.Movement{
			ID:            uuid.New().String(),
			AccountID:     uuid.New().String(),
			InstitutionID: "iid",
			Type:          domain.Income,
			Category:      domain.Education,
			Amount:        150.00,
			Date:          time.Now(),
		}

		mockRepo.On("CreateMovement", ctx, movement).Return(nil).Once()

		err := u.CreateMovement(ctx, movement)
		c.NoError(err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("nil movement", func(t *testing.T) {
		err := u.CreateMovement(ctx, nil)
		c.ErrorContains(err, "movement cannot be nil")
	})

	t.Run("missing accountID", func(t *testing.T) {
		m := &domain.Movement{
			InstitutionID: "iid",
			Type:          domain.Income,
			Amount:        100,
			Date:          time.Now(),
		}
		err := u.CreateMovement(ctx, m)
		c.ErrorContains(err, "account ID is required")
	})

	t.Run("invalid amount", func(t *testing.T) {
		m := &domain.Movement{
			AccountID:     uuid.New().String(),
			InstitutionID: "iid",
			Type:          domain.Income,
			Amount:        0,
			Date:          time.Now(),
		}
		err := u.CreateMovement(ctx, m)
		c.ErrorContains(err, "amount must be greater than zero")
	})

	t.Run("invalid type", func(t *testing.T) {
		m := &domain.Movement{
			AccountID:     uuid.New().String(),
			InstitutionID: "iid",
			Type:          "other",
			Amount:        100,
			Date:          time.Now(),
		}
		err := u.CreateMovement(ctx, m)
		c.Error(err)
	})

	t.Run("date in the future", func(t *testing.T) {
		m := &domain.Movement{
			AccountID:     uuid.New().String(),
			InstitutionID: "iid",
			Type:          domain.Income,
			Amount:        100,
			Category:      domain.Education,
			Date:          time.Now().Add(24 * time.Hour),
		}
		err := u.CreateMovement(ctx, m)
		c.ErrorContains(err, "movement date cannot be in the future")
	})

	t.Run("repository error", func(t *testing.T) {
		m := &domain.Movement{
			AccountID:     uuid.New().String(),
			InstitutionID: "iid",
			Type:          domain.Income,
			Amount:        100,
			Category:      domain.Education,
			Date:          time.Now(),
		}

		expectedErr := errors.New("repo failure")
		mockRepo.On("CreateMovement", ctx, m).Return(expectedErr).Once()

		err := u.CreateMovement(ctx, m)
		c.ErrorIs(err, expectedErr)

		mockRepo.AssertExpectations(t)
	})
}
func TestCreateMovementWithRepositoryError(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)
	usecase := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()

	testMovement := &domain.Movement{
		AccountID:     uuid.New().String(),
		InstitutionID: "iid",
		Category:      domain.Debt,
		Type:          domain.Expense,
		Amount:        150.00,
	}

	dbErr := errors.New("db error")

	mockRepo.On("CreateMovement", ctx, testMovement).Return(dbErr).Once()

	err := usecase.CreateMovement(ctx, testMovement)
	c.ErrorIs(err, dbErr)

	mockRepo.AssertExpectations(t)
}

func TestGetMovementByID(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)
	usecase := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()
	testID := uuid.New().String()
	expectedMovement := &domain.Movement{ID: testID, AccountID: "acc1"}

	mockRepo.On("GetMovementByID", ctx, testID, "acc1").Return(expectedMovement, nil).Once()
	foundMovement, err := usecase.GetMovementByID(ctx, testID, "acc1")
	c.NoError(err)
	c.NotNil(foundMovement)
	c.Equal(expectedMovement.ID, foundMovement.ID)
	c.Equal(expectedMovement.AccountID, foundMovement.AccountID)
	mockRepo.AssertExpectations(t)
}

func TestGetMovementByIDWithRepositoryError(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)
	usecase := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()
	testID := uuid.New().String()

	mockRepo.On("GetMovementByID", ctx, testID, "acc1").Return(nil, errors.New("db error")).Once()
	foundMovement, err := usecase.GetMovementByID(ctx, testID, "acc1")
	c.Error(err)
	c.Nil(foundMovement)
	mockRepo.AssertExpectations(t)
}

func TestGetMovementsByAccountID(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)
	usecase := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()

	testAccountID := uuid.New().String()

	expectedMovements := []*domain.Movement{
		{ID: uuid.New().String(), AccountID: testAccountID, Amount: 100},
		{ID: uuid.New().String(), AccountID: testAccountID, Amount: 200},
	}

	limit := 10
	offset := 0

	mockRepo.On("GetMovementsByAccountID", ctx, testAccountID, limit, offset).
		Return(expectedMovements, nil).Once()

	mockRepo.On("GetTotalMovementsByAccountID", ctx, testAccountID).
		Return(len(expectedMovements), nil)

	foundMovements, err := usecase.GetPaginatedMovementsByAccountID(ctx, testAccountID, limit, offset)
	c.NoError(err)
	c.NotNil(foundMovements)
	c.Equal(len(expectedMovements), len(foundMovements.Movements))

	mockRepo.AssertExpectations(t)
}

func TestGetMovementsByAccountIDWithRepositoryError(t *testing.T) {
	c := require.New(t)
	mockRepo := new(repository.MockMovementRepository)
	usecase := NewMovementUsecase(context.Background(), mockRepo)
	ctx := context.Background()
	testAccountID := uuid.New().String()

	limit := 10
	offset := 0

	mockRepo.On("GetMovementsByAccountID", ctx, testAccountID, limit, offset).
		Return(nil, errors.New("db error")).Once()

	mockRepo.On("GetTotalMovementsByAccountID", ctx, testAccountID).
		Return(1, nil)

	foundMovements, err := usecase.GetPaginatedMovementsByAccountID(ctx, testAccountID, limit, offset)
	c.Error(err)
	c.Nil(foundMovements)

	mockRepo.AssertExpectations(t)
}
