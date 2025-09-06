package google

import (
	"errors"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GoogleDeleteWath() gin.HandlerFunc {
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

		gmailService, err := gClient.GmailService(c)
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

		err = gmailService.DeleteWatch()
		if errors.Is(err, googleapi.ErrMissingHistoryID) {
			log.Info(loggerModels.LogProperties{
				Event: "missing_history_id",
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "missing historyID",
			})

			return
		}

		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "delete_watch_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "watch deleted succefully",
		})
	}
}
