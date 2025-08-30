package google

import (
	"fmt"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, err := logger.GetLogger(c, "google-login")
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: fmt.Sprintf("logger not init: %s", err.Error()),
			})

			return
		}

		err = gClient.SaveTokenAndInitServices(c, c.Query("code"))
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "save_token_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		projectID := "transaction-tracker-2473"
		topicName := fmt.Sprintf("projects/%s/topics/gmail-notifications", projectID)

		gmailService, err := gClient.GmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_gmail_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		historyID, expirationTime, err := gmailService.CreateWatch(c, topicName)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "create_watch_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: fmt.Sprintf("HistoryID: %d ExpirationTime: %d", historyID, expirationTime),
		})
	}
}
