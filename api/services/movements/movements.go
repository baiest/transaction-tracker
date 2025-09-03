package services

import (
	"context"
	"fmt"
	"time"
	"transaction-tracker/api/repositories"
	"transaction-tracker/database/mongo/schemas"
)

type MovementsService struct {
	repo repositories.IMovementsRepository
}

func NewMovementsService(ctx context.Context) (*MovementsService, error) {
	repo, err := repositories.NewMovementsRepository(ctx)
	if err != nil {
		return nil, err
	}

	return &MovementsService{repo: repo}, nil
}

func (s *MovementsService) CreateMovement(ctx context.Context, movement *schemas.Movement) error {
	fmt.Println("Creating movement:", movement)
	return s.repo.SaveMovement(ctx, movement)
}

func (s *MovementsService) GetMovements(ctx context.Context, startDate time.Time, finishDate time.Time) ([]*schemas.Movement, error) {
	return s.repo.GetMovements(ctx, startDate, finishDate)
}
