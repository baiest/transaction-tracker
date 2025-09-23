package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/repository"
)

type movementUsecase struct {
	movementRepo repository.MovementRepository
}

// NewMovementUsecase is the constructor for the use case implementation.
// It receives a repository interface as a dependency.
func NewMovementUsecase(repo repository.MovementRepository) MovementUsecase {
	return &movementUsecase{
		movementRepo: repo,
	}
}

// CreateMovement contains the business logic for creating a movement.
func (u *movementUsecase) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	if movement == nil {
		return errors.New("movement cannot be nil")
	}

	if movement.AccountID == "" {
		return errors.New("account ID is required")
	}

	if movement.InstitutionID == "" {
		return errors.New("institution ID is required")
	}

	if movement.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	_, err := domain.ParseMovementType(string(movement.Type))
	if err != nil {
		return err
	}

	if movement.Date.After(time.Now()) {
		return errors.New("movement date cannot be in the future")
	}

	return u.movementRepo.CreateMovement(ctx, movement)
}

// GetMovementByID is a sample method to get a movement.
func (u *movementUsecase) GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error) {
	return u.movementRepo.GetMovementByID(ctx, id, accountID)
}

// GetMovementsByUserID is a sample method to get a user's movements.
func (u *movementUsecase) GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) (*domain.PaginatedMovements, error) {
	fmt.Println(limit, offset)
	if limit <= 0 {
		limit = 10
	}

	if limit > 20 {
		limit = 20
	}

	offset -= 1
	if offset < 0 {
		offset = 0
	}

	totalRecords, err := u.movementRepo.GetTotalMovementsByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if limit > 0 {
		totalPages = (totalRecords + limit - 1) / limit
	}

	movements, err := u.movementRepo.GetMovementsByAccountID(ctx, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &domain.PaginatedMovements{
		Movements:    movements,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
		Limit:        limit,
		Offset:       offset,
		CurrentPage:  (offset / limit) + 1,
	}, nil
}
