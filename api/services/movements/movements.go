package movements

import (
	"context"
	"errors"
	"time"
	"transaction-tracker/api/models"
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
	return s.repo.SaveMovement(ctx, movement)
}

func (s *MovementsService) GetMovements(ctx context.Context, page int64) ([]*schemas.Movement, int64, error) {
	return s.repo.GetMovements(ctx, page)
}

func (s *MovementsService) GetMovementsByMonth(ctx context.Context, year int, month int) (*models.MovementByMonth, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	finishDate := time.Date(year, time.Month(month+1), 0, 23, 59, 59, 999999999, time.UTC)

	movements, err := s.repo.GetMovementsByDateRange(ctx, startDate, finishDate)
	if err != nil {
		return nil, err
	}

	movementsByMonth := &models.MovementByMonth{}
	movementsByDay := make([]*models.MovementIncomeOutcomeByDay, 31)

	for _, m := range movements {
		m.Date = m.Date.In(time.Local)
		day := m.Date.Day() - 1

		if movementsByDay[day] == nil {
			movementsByDay[day] = &models.MovementIncomeOutcomeByDay{Day: day + 1}
		}

		if m.IsNegative {
			movementsByDay[day].Outcome += m.Value
			movementsByMonth.TotalOutcome += m.Value
		} else {
			movementsByDay[day].Income += m.Value
			movementsByMonth.TotalIncome += m.Value
		}
	}

	movementsByMonth.Days = []*models.MovementIncomeOutcomeByDay{}

	for _, v := range movementsByDay {
		if v != nil {
			movementsByMonth.Days = append(movementsByMonth.Days, v)
		}
	}

	movementsByMonth.Balance = movementsByMonth.TotalIncome - movementsByMonth.TotalOutcome
	movementsByMonth.Year = year

	return movementsByMonth, nil
}

func (s *MovementsService) GetMovementsByYear(ctx context.Context, year int) (*models.MovementByYear, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	finishDate := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	movements, err := s.repo.GetMovementsByDateRange(ctx, startDate, finishDate)
	if err != nil {
		return nil, err
	}

	movementsByYear := &models.MovementByYear{}
	movementsByMonth := make([]*models.MovementIncomeOutcome, 12)

	for _, m := range movements {
		month := m.Date.Month() - 1
		if movementsByMonth[month] == nil {
			movementsByMonth[month] = &models.MovementIncomeOutcome{}
		}

		if m.IsNegative {
			movementsByMonth[month].Outcome += m.Value
			movementsByYear.TotalOutcome += m.Value
		} else {
			movementsByMonth[month].Income += m.Value
			movementsByYear.TotalIncome += m.Value
		}
	}

	movementsByYear.Months = []*models.MovementIncomeOutcome{}

	for _, v := range movementsByMonth {
		if v != nil {
			movementsByYear.Months = append(movementsByYear.Months, v)
		}
	}

	movementsByYear.Balance = movementsByYear.TotalIncome - movementsByYear.TotalOutcome

	return movementsByYear, nil
}

func (s *MovementsService) DeleteMovement(ctx context.Context, movementID string) error {
	if movementID == "" {
		return errors.New("movementID is required")
	}

	return s.repo.DeleteMovement(ctx, movementID)
}
