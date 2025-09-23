package handler

import (
	"errors"
	"transaction-tracker/api/models"
	"transaction-tracker/internal/messages/repository"
	"transaction-tracker/internal/messages/usecase"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

// MessageHandler handles HTTP requests for the messages domain.
type MessageHandler struct {
	messageUsecase usecase.MessageUsecase
}

// NewMessageHandler creates a new instance of MessageHandler.
func NewMessageHandler(ucm usecase.MessageUsecase) *MessageHandler {
	return &MessageHandler{
		messageUsecase: ucm,
	}
}

// GetMessageByID handles the GET /messages/:id request.
func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	id := c.Param("id")
	if id == "" {
		models.NewResponseInvalidRequest(c, models.Response{Message: "message id is required"})
		return
	}

	message, err := h.messageUsecase.GetMessageByIDAndAccountID(c.Request.Context(), id, account.ID)
	if err != nil {
		if errors.Is(err, repository.ErrMessageNotFound) {
			models.NewResponseNotFound(c, models.Response{Message: "message not found"})
			return
		}

		log.Error(loggerModels.LogProperties{
			Event: "get_message_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseOK(c, models.Response{
		Data: models.ToMessageResponse(message),
	})
}

// ProcessMessage handles the POST /messages request.
func (h *MessageHandler) ProcessMessage(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	var req models.MessageRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "invalid_request_body",
			Error: err,
		})

		models.NewResponseInvalidRequest(c, models.Response{Message: "invalid form data"})
		return
	}

	message, err := h.messageUsecase.Process(c.Request.Context(), "", req.ExternalID, account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "process_message_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)
		return
	}

	models.NewResponseOK(c, models.Response{
		Data: models.ToMessageResponse(message),
	})
}
