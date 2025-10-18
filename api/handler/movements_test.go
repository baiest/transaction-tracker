package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"transaction-tracker/api/models"
	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/usecase"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetMovements_Success(t *testing.T) {
	c := require.New(t)

	movements := []*domain.Movement{
		{
			ID:        "1",
			AccountID: "test-account-id",
			Amount:    100.0,
			Date:      time.Now(),
		},
	}

	movementsMock := new(usecase.MockMovementUsecase)
	movementsMock.On("GetPaginatedMovementsByAccountID", mock.Anything, "accountID", 5, 2).Return(&domain.PaginatedMovements{Movements: movements, CurrentPage: 1}, nil)

	testHandler := NewMovementHandler(movementsMock)

	ginContext, w := setupTestContext(http.MethodGet, "/movements?page=2&limit=5", nil)

	testHandler.GetMovements(ginContext)

	c.Equal(http.StatusOK, w.Code)

	var response *models.MovementsListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	c.NoError(err)

	c.Len(response.Movements, 1)
	c.Equal(int64(1), response.Page)
}

func TestGetMovements_UsecaseError(t *testing.T) {
	c := require.New(t)

	mockUsecase := new(usecase.MockMovementUsecase)
	mockUsecase.On("GetPaginatedMovementsByAccountID", mock.Anything, "accountID", 5, 2).Return(nil, errors.New("database connection failed"))

	testHandler := NewMovementHandler(mockUsecase)

	ginContext, w := setupTestContext(http.MethodGet, "/movements?page=2&limit=5", nil)

	testHandler.GetMovements(ginContext)

	c.Equal(http.StatusInternalServerError, w.Code)
}

func TestCreateMovement_UsecaseError(t *testing.T) {
	c := require.New(t)

	mockUsecase := new(usecase.MockMovementUsecase)
	mockUsecase.On("CreateMovement", mock.Anything, mock.Anything).Return(errors.New("database connection failed"))

	testHandler := NewMovementHandler(mockUsecase)

	body := strings.NewReader("institution_id=id&category=food&type=income&amount=1500&date=2025-09-20T10:30:00Z&description=test+movement")
	ginContext, w := setupTestContext(http.MethodPost, "/movements", body)

	ginContext.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler.CreateMovement(ginContext)

	c.Equal(http.StatusInternalServerError, w.Code)
}

func TestCreateMovement_InvalidRequestBody(t *testing.T) {
	c := require.New(t)

	mockUsecase := new(usecase.MockMovementUsecase)
	mockUsecase.On("CreateMovement", mock.Anything, mock.Anything).Return(nil)

	testHandler := NewMovementHandler(mockUsecase)

	body := strings.NewReader("type=income&amount=abc&date=not-a-date")
	ginContext, w := setupTestContext(http.MethodPost, "/movements", body)
	ginContext.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler.CreateMovement(ginContext)

	c.Equal(http.StatusBadRequest, w.Code)
	c.Contains(w.Body.String(), "invalid form data")
}

func TestCreateMovement_Success(t *testing.T) {
	c := require.New(t)

	mockUsecase := new(usecase.MockMovementUsecase)
	mockUsecase.On("CreateMovement", mock.Anything, mock.Anything).Return(nil)

	testHandler := NewMovementHandler(mockUsecase)

	body := strings.NewReader(
		"type=income&institution_id=inst-1&category=food&description=Salary&amount=1500&date=2025-09-20T10:17:00Z",
	)

	ginContext, w := setupTestContext(http.MethodPost, "/movements", body)
	ginContext.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler.CreateMovement(ginContext)

	c.Equal(http.StatusCreated, w.Code)
}
