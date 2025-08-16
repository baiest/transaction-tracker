package google

import (
	"fmt"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := gClient.SaveTokenAndInitServices(c, c.Query("code"))
		if err != nil {
			models.NewResponseInternalServerError(c)

			return
		}

		projectID := "transaction-tracker-2473"
		topicName := fmt.Sprintf("projects/%s/topics/gmail-notifications", projectID)

		gmailService, err := gClient.GmailService()
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		historyID, expirationTime, err := gmailService.CreateWatch(c, topicName)
		if err != nil {
			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: fmt.Sprintf("HistoryID: %d ExpirationTime: %d", historyID, expirationTime),
		})
	}
}
