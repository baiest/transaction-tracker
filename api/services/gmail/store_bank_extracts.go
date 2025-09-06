package gmail

import (
	"transaction-tracker/api/models"
	"transaction-tracker/database/mongo/schemas"

	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func StoreBankExtracts() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		gmailService, err := NewGmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_google_service_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)
			return
		}

		extracts, err := gmailService.GetExtractMessages(c, "davivienda")
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_extract_messages_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		messages := []*schemas.Message{}
		messagesChan := make(chan bool, len(extracts.Messages))

		for _, msg := range extracts.Messages {
			go func() {
				defer func() {
					messagesChan <- true
				}()

				message, err := gmailService.ProcessMessage(c, msg.Id, "")
				if err != nil {
					log.Error(loggerModels.LogProperties{
						Event: "process_message_failed",
						Error: err,
						AdditionalParams: []loggerModels.Properties{
							logger.MapToProperties(map[string]string{
								"message_id": msg.Id,
							}),
						},
					})

					return
				}

				messages = append(messages, message)
			}()
		}

		for range extracts.Messages {
			<-messagesChan
		}

		close(messagesChan)

		models.NewResponseOK(c, models.Response{
			Data: messages,
		})
	}
}
