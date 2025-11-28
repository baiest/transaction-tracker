package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/repository"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"
	"transaction-tracker/shared"
)

type ClassifyResponse struct {
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
}

var (
	ErrMovementNotFound      = errors.New("movement not found")
	ErrMustBeGreaterThanZero = errors.New("amount must be greater than zero")
)

type movementUsecase struct {
	movementRepo repository.MovementRepository
	log          *loggerModels.Logger
}

// NewMovementUsecase is the constructor for the use case implementation.
// It receives a repository interface as a dependency.
func NewMovementUsecase(ctx context.Context, repo repository.MovementRepository) MovementUsecase {
	log, _ := logger.GetLogger(ctx, "movements-usecase")

	return &movementUsecase{
		movementRepo: repo,
		log:          log,
	}
}

// CreateMovement contains the business logic for creating a movement.
func (u *movementUsecase) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	if movement == nil {
		return errors.New("movement cannot be nil")
	}

	if movement.InstitutionID == "" {
		return errors.New("institution ID is required")
	}

	if movement.AccountID == "" {
		return errors.New("account ID is required")
	}

	if movement.Amount <= 0 {
		return ErrMustBeGreaterThanZero
	}

	_, err := domain.ParseMovementCategory(string(movement.Category))
	if err != nil {
		return err
	}

	_, err = domain.ParseMovementType(string(movement.Type))
	if err != nil {
		return err
	}

	if movement.Date.After(time.Now()) {
		return errors.New("movement date cannot be in the future")
	}

	movement.Category = domain.Unknown

	if movement.Description != "" {
		err = u.ClassifyCategoryFromDescription(ctx, movement)
		if err != nil {
			u.log.Error(loggerModels.LogProperties{
				Event: "error_classifying_category",
				Error: err,
				AdditionalParams: []loggerModels.Properties{
					movement,
				},
			})
		}
	}

	return u.movementRepo.CreateMovement(ctx, movement)
}

// GetMovementByID is a sample method to get a movement.
func (u *movementUsecase) GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error) {
	movement, err := u.movementRepo.GetMovementByID(ctx, id, accountID)
	if err != nil {
		if errors.Is(err, repository.ErrMovementNotFound) {
			return nil, ErrMovementNotFound
		}

		return nil, err
	}

	return movement, nil
}

// GetMovementsByUserID is a sample method to get a user's movements.
func (u *movementUsecase) GetPaginatedMovementsByAccountID(ctx context.Context, accountID string, institutionIDs []string, limit int, offset int) (*domain.PaginatedMovements, error) {
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

	movements, err := u.movementRepo.GetMovementsByAccountID(ctx, accountID, institutionIDs, limit, offset)
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

func (u *movementUsecase) DeleteMovement(ctx context.Context, id string, accountID string) error {
	_, err := u.GetMovementByID(ctx, id, accountID)
	if err != nil {
		return err
	}

	return u.movementRepo.Delete(ctx, id, accountID)
}

func (u *movementUsecase) GetMovementsByYear(ctx context.Context, accountID string, institutionIDs []string, year int) ([]*domain.Movement, error) {
	movements, err := u.movementRepo.GetMovementsByAccountID(ctx, accountID, institutionIDs, 1000, 0)
	if err != nil {
		return nil, err
	}

	filteredMovements := []*domain.Movement{}
	for _, m := range movements {
		if m.Date.Year() == year {
			filteredMovements = append(filteredMovements, m)
		}
	}

	return filteredMovements, nil
}

func (u *movementUsecase) GetMovementsByMonth(ctx context.Context, accountID string, institutionIDs []string, year int, month int) ([]*domain.Movement, error) {
	movements, err := u.movementRepo.GetMovementsByAccountID(ctx, accountID, institutionIDs, 1000, 0)
	if err != nil {
		return nil, err
	}

	filteredMovements := []*domain.Movement{}
	for _, m := range movements {
		if m.Date.Year() == year && int(m.Date.Month()) == month {
			filteredMovements = append(filteredMovements, m)
		}
	}

	return filteredMovements, nil
}

func (u *movementUsecase) DeleteMovementsByExtractID(ctx context.Context, extractID string) error {
	return u.movementRepo.DeleteMovementsByExtractID(ctx, extractID)
}

func (u *movementUsecase) ClassifyCategoryFromDescription(ctx context.Context, movement *domain.Movement) error {
	payload := map[string]interface{}{
		"description": movement.Description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", os.Getenv("CLASSIFY_CATEGORY_URL"), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := shared.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error classifying category: status %d", res.StatusCode)
	}

	cr := ClassifyResponse{}
	err = json.NewDecoder(res.Body).Decode(&cr)
	if err != nil {
		return err
	}

	u.log.Info(loggerModels.LogProperties{
		Event: "category_classified",
		AdditionalParams: []loggerModels.Properties{
			logger.MapToProperties(map[string]string{
				"confidence": string(fmt.Sprintf("%.2f", cr.Confidence)),
				"category":   cr.Category,
			}),
		},
	})

	movement.Category = domain.MovementCategory(cr.Category)

	return nil
}
