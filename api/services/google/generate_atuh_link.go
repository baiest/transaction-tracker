package google

import (
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

func GoogleGenerateAuthLink(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		models.NewResponseOK(c, models.Response{
			Message: gClient.GetAuthURL(),
		})
	}
}
