package handler

import (
	"transaction-tracker/api/models"
	"transaction-tracker/internal/extracts/usecase"
	messagesUsecase "transaction-tracker/internal/messages/usecase"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

type ExtractsHandler struct {
	extractsUsecase usecase.ExtractsUsecase
	msgUsecase      messagesUsecase.MessageUsecase
}

// NewAccountHandler creates a new instance of AccountsHandler.
func NewExtractsHandler(uce usecase.ExtractsUsecase, ucm messagesUsecase.MessageUsecase) *ExtractsHandler {
	return &ExtractsHandler{
		extractsUsecase: uce,
		msgUsecase:      ucm,
	}
}

func (h *ExtractsHandler) GetAllExtracts(c *gin.Context) {
	log, account, err := getContextDependencies(c)
	if err != nil {
		return
	}

	messageIDs, err := h.extractsUsecase.GetExtractMessages(c.Request.Context(), "davivienda", account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_extract_messages_failed",
			Error: err,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	messagesProcessed := []*models.MessageResponse{}

	for _, id := range messageIDs {
		message, err := h.msgUsecase.Process(c.Request.Context(), "", id, account)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "process_message_failed",
				Error: err,
			})
		}

		messagesProcessed = append(messagesProcessed, models.ToMessageResponse(message))
	}

	models.NewResponseOK(c, models.Response{
		Data: messagesProcessed,
	})
}
