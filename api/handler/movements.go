package handler

import (
	"errors"
	"strconv"
	"time"
	"transaction-tracker/api/models"
	"transaction-tracker/internal/movements/domain"
	"transaction-tracker/internal/movements/usecase"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MovementHandler handles HTTP requests for the movements domain.
type MovementHandler struct {
	movementsUsecase usecase.MovementUsecase
}

// NewMovementHandler creates a new instance of MovementHandler.
func NewMovementHandler(ucm usecase.MovementUsecase) *MovementHandler {
	return &MovementHandler{
		movementsUsecase: ucm,
	}
}

// GetMovements handles the GET /movements request.
func (h *MovementHandler) GetMovements(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}

	movements, err := h.movementsUsecase.GetPaginatedMovementsByAccountID(c.Request.Context(), account.ID, int(limit), int(page))
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_movements_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseOK(c, models.Response{
		Data: &models.MovementsListResponse{
			TotalPages: int64(movements.TotalPages),
			Page:       int64(movements.CurrentPage),
			Movements:  models.ToMovementResponses(movements.Movements),
		},
	})
}

// GetMovementByID handles the GET /movements request.
func (h *MovementHandler) GetMovementByID(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewResponseInvalidRequest(c, models.Response{Message: "movement id is required"})
		return
	}

	movement, err := h.movementsUsecase.GetMovementByID(c.Request.Context(), id, account.ID)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_movements_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseOK(c, models.Response{
		Data: models.ToMovementResponse(movement),
	})
}

// CreateMovement handles the POST /movements request.
func (h *MovementHandler) CreateMovement(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	var req models.CreateMovementRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "invalid_request_body",
			Error: err,
		})

		errorMessage := "invalid form data"

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			firstError := validationErrors[0]
			field := firstError.Field() // Nombre del campo en la struct (ej: "Amount")
			tag := firstError.Tag()     // Etiqueta de validación (ej: "required")

			if tag == "required" {
				errorMessage = field + " is required"
			} else {
				errorMessage = field + " has invalid value"
			}
		}

		models.NewResponseInvalidRequest(c, models.Response{Message: errorMessage})
		return
	}

	req.AccountID = account.ID

	movement := models.ToDomainMovement(req)

	if err := h.movementsUsecase.CreateMovement(c.Request.Context(), movement); err != nil {
		if errors.Is(err, domain.ErrInvalidMovementType) || errors.Is(err, domain.ErrInvalidMovementCategory) {
			models.NewResponseInvalidRequest(c, models.Response{Message: err.Error()})
			return
		}

		log.Error(loggerModels.LogProperties{
			Event: "create_movement_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseCreated(c, models.Response{
		Data: models.ToMovementResponse(movement),
	})
}

// DeleteMovement handles DELETE /movements/:id request
func (h *MovementHandler) DeleteMovement(c *gin.Context) {
	log, accconst, err := getContextDependencies(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewResponseInvalidRequest(c, models.Response{Message: "movement id is required"})
		return
	}

	err = h.movementsUsecase.DeleteMovement(c.Request.Context(), id, accconst.ID)
	if err != nil {
		if errors.Is(err, usecase.ErrMovementNotFound) {
			models.NewResponseNotFound(c, models.Response{Message: "movement not found"})
			return
		}

		log.Error(loggerModels.LogProperties{
			Event: "delete_movement_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseOK(c, models.Response{
		Message: "movement deleted successfully",
	})
}

func (h *MovementHandler) GetMovementsByYear(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "invalid_year",
			Error: err,
		})

		models.NewResponseInvalidRequest(c, models.Response{
			Message: "invalid year",
		})

		return
	}

	movements, err := h.movementsUsecase.GetMovementsByYear(c.Request.Context(), account.ID, year)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_movements_by_year_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	totalIncome := 0.0
	totalOutcome := 0.0
	monthsMap := map[time.Month]models.MovementIncomeOutcomeByMonth{}

	for _, m := range movements {
		if m.Type == domain.Income {
			totalIncome += m.Amount
		}

		if m.Type == domain.Expense {
			totalOutcome += m.Amount
		}

		mio := models.MovementIncomeOutcomeByMonth{
			Month: m.Date.Month(),
			MovementIncomeOutcome: models.MovementIncomeOutcome{
				Income:  0.0,
				Outcome: 0.0,
			},
		}

		_, ok := monthsMap[mio.Month]
		if ok {
			mio.Income = monthsMap[mio.Month].Income
			mio.Outcome = monthsMap[mio.Month].Outcome
		}

		if m.Type == domain.Income {
			mio.Income += m.Amount
		}

		if m.Type == domain.Expense {
			mio.Outcome += m.Amount
		}

		monthsMap[mio.Month] = mio
	}

	months := []models.MovementIncomeOutcomeByMonth{}
	for _, m := range monthsMap {
		months = append(months, m)
	}

	models.NewResponseOK(c, models.Response{
		Data: models.MovementByYear{
			TotalIncome:  float64(totalIncome),
			TotalExpense: float64(totalOutcome),
			Balance:      float64(totalIncome) - float64(totalOutcome),
			Months:       months,
		},
	})
}

func (h *MovementHandler) GetMovementsByMonth(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	yearStr := c.Param("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "invalid_year",
			Error: err,
		})

		models.NewResponseInvalidRequest(c, models.Response{
			Message: "invalid year",
		})

		return
	}

	monthStr := c.Param("month")
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "invalid_month",
			Error: err,
		})

		models.NewResponseInvalidRequest(c, models.Response{
			Message: "invalid month",
		})

		return
	}

	movements, err := h.movementsUsecase.GetMovementsByMonth(c.Request.Context(), account.ID, year, month)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_movements_by_month_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	totalIncome := 0.0
	totalOutcome := 0.0
	daysMap := map[int]models.MovementIncomeOutcomeByDay{}

	for _, m := range movements {
		if m.Type == domain.Income {
			totalIncome += m.Amount
		}

		if m.Type == domain.Expense {
			totalOutcome += m.Amount
		}

		mio := models.MovementIncomeOutcomeByDay{
			Day: m.Date.Day(),
			MovementIncomeOutcome: models.MovementIncomeOutcome{
				Income:  0.0,
				Outcome: 0.0,
			},
		}

		_, ok := daysMap[mio.Day]
		if ok {
			mio.Income = daysMap[mio.Day].Income
			mio.Outcome = daysMap[mio.Day].Outcome
		}

		if m.Type == domain.Income {
			mio.Income += m.Amount
		}

		if m.Type == domain.Expense {
			mio.Outcome += m.Amount
		}

		daysMap[mio.Day] = mio
	}

	days := []models.MovementIncomeOutcomeByDay{}
	for _, m := range daysMap {
		days = append(days, m)
	}

	models.NewResponseOK(c, models.Response{
		Data: models.MovementByMonth{
			TotalIncome:  float64(totalIncome),
			TotalOutcome: float64(totalOutcome),
			Year:         year,
			Balance:      float64(totalIncome) - float64(totalOutcome),
			Days:         days,
		},
	})
}
