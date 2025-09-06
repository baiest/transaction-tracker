package gmail

import (
	"transaction-tracker/api/models"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		messageID := c.Param("messageID")
		if messageID == "" {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "missing message id",
			})

			return
		}

		gmailService, err := NewGmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_gmail_service_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		message, err := gmailService.GetMessage(c, messageID)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_message_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: message,
		})
	}
}
