package google

import (
	"strconv"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GetEmailByHistoryID() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		gClient, err := googleapi.NewGoogleClient(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_google_client_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)
			return
		}

		historyID, err := strconv.ParseUint(c.Param("historyID"), 10, 64)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "invalid_history_id",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid historyID",
			})

			return
		}

		gmailClient, err := gClient.GmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_gmail_service_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		historyListCall := gmailClient.Client.Users.History.List("me").StartHistoryId(historyID)

		historyList, err := historyListCall.Do()
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_history_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		if len(historyList.History) == 0 {
			log.Error(loggerModels.LogProperties{
				Event: "no_emails_found",
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "no emails found for given historyID",
			})

			return
		}

		var messages []map[string]interface{}
		for _, h := range historyList.History {
			for _, m := range h.Messages {
				msg, err := gmailClient.Client.Users.Messages.Get("me", m.Id).Format("full").Do()
				if err != nil {
					continue
				}

				messages = append(messages, map[string]interface{}{
					"id":      msg.Id,
					"snippet": msg.Snippet,
					"payload": msg.Payload,
				})
			}
		}

		models.NewResponseOK(c, models.Response{
			Data: map[string]interface{}{
				"historyID": historyID,
				"messages":  messages,
			},
		})
	}
}
