package handler

import (
	"strconv"
	"transaction-tracker/api/models"
	"transaction-tracker/internal/movements/usecase"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
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

	limitStr := c.DefaultQuery("limit", "1")
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 1
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

		models.NewResponseInvalidRequest(c, models.Response{Message: "invalid form data"})
		return
	}

	req.AccountID = account.ID

	movement := models.ToDomainMovement(req)

	if err := h.movementsUsecase.CreateMovement(c.Request.Context(), movement); err != nil {
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
