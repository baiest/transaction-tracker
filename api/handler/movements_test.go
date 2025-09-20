package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"transaction-tracker/api/models"
	"transaction-tracker/api/services/accounts"
	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/usecase"
	"transaction-tracker/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupTestContext(method, target string, body io.Reader) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest(method, target, body)
	if err != nil {
	}
	c.Request = req

	mockLogger, _ := logger.GetLogger(c, "test")
	c.Set("logger", mockLogger)

	c.Set("account", &accounts.Account{ID: "accountID"})

	return c, w
}

func TestGetMovements_Success(t *testing.T) {
	c := require.New(t)

	movementsMock := []*domain.Movement{
		{
			ID:        "1",
			AccountID: "test-account-id",
			Amount:    100.0,
			Date:      time.Now(),
		},
	}

	mockUsecase := &usecase.MockMovementUsecase{
		GetPaginatedMovementsByAccountIDFunc: func(ctx context.Context, accountID string, limit, offset int) (*domain.PaginatedMovements, error) {
			return &domain.PaginatedMovements{Movements: movementsMock, CurrentPage: 1}, nil
		},
	}
	testHandler := NewMovementHandler(mockUsecase)

	ginContext, w := setupTestContext(http.MethodGet, "/movements?page=2&limit=5", nil)

	testHandler.GetMovements(ginContext)

	c.Equal(http.StatusOK, w.Code)

	var response *models.MovementsListResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	c.NoError(err)

	fmt.Println(response)

	c.Len(response.Movements, 1)
	c.Equal(int64(1), response.Page)
}

func TestGetMovements_UsecaseError(t *testing.T) {
	c := require.New(t)

	mockUsecase := &usecase.MockMovementUsecase{
		GetPaginatedMovementsByAccountIDFunc: func(ctx context.Context, accountID string, limit, offset int) (*domain.PaginatedMovements, error) {
			return nil, errors.New("database connection failed")
		},
	}
	testHandler := NewMovementHandler(mockUsecase)

	ginContext, w := setupTestContext(http.MethodGet, "/movements?page=2&limit=5", nil)

	testHandler.GetMovements(ginContext)

	c.Equal(http.StatusInternalServerError, w.Code)
}

func TestCreateMovement_UsecaseError(t *testing.T) {
	c := require.New(t)

	mockUsecase := &usecase.MockMovementUsecase{
		CreateMovementFunc: func(ctx context.Context, movement *domain.Movement) error {
			return errors.New("db error")
		},
	}

	testHandler := NewMovementHandler(mockUsecase)

	body := strings.NewReader("institution_id=iid&type=income&amount=1500&date=2025-09-20T10:30:00Z&description=test+movement")
	ginContext, w := setupTestContext(http.MethodPost, "/movements", body)

	ginContext.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler.CreateMovement(ginContext)

	c.Equal(http.StatusInternalServerError, w.Code)
}

func TestCreateMovement_InvalidRequestBody(t *testing.T) {
	c := require.New(t)

	mockUsecase := &usecase.MockMovementUsecase{
		CreateMovementFunc: func(ctx context.Context, movement *domain.Movement) error {
			t.Fatalf("CreateMovement should not be called in case of a binding error")
			return nil
		},
	}
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

	called := false
	mockUsecase := &usecase.MockMovementUsecase{
		CreateMovementFunc: func(ctx context.Context, movement *domain.Movement) error {
			called = true

			c.Equal(domain.Income, movement.Type)
			c.Equal("accountID", movement.AccountID)
			c.Equal("inst-1", movement.InstitutionID)
			c.Equal("Salary", movement.Description)
			c.Equal(float64(1500), movement.Amount)
			c.Equal("2025-09-20", movement.Date.Format("2006-01-02"))

			return nil
		},
	}
	testHandler := NewMovementHandler(mockUsecase)

	body := strings.NewReader(
		"type=income&institution_id=inst-1&description=Salary&amount=1500&date=2025-09-20T10:17:00Z",
	)

	ginContext, w := setupTestContext(http.MethodPost, "/movements", body)
	ginContext.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	testHandler.CreateMovement(ginContext)

	c.True(called, "usecase.CreateMovement should be called")
	c.Equal(http.StatusCreated, w.Code)
}
