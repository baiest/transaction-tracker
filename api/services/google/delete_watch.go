package google

import (
	"errors"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

func GoogleDeleteWath(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		gmailService, err := gClient.GmailService()
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		err = gmailService.DeleteWatch()
		if errors.Is(err, googleapi.ErrMissingHistoryID) {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "missing historyID",
			})

			return
		}

		if err != nil {
			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "watch deleted succefully",
		})
	}
}
