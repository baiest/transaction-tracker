package google

import (
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func GoogleGenerateAuthLink() gin.HandlerFunc {
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

		models.NewResponseOK(c, models.Response{
			Message: gClient.GetAuthURL(),
		})
	}
}
