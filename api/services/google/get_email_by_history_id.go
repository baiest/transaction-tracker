package google

import (
	"strconv"
	"transaction-tracker/api/models"
	"transaction-tracker/googleapi"

	"github.com/gin-gonic/gin"
)

func GetEmailByHistoryID(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		if email == "" {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "email is required as a query parameter",
			})
			return
		}

		historyID, err := strconv.ParseUint(c.Param("historyID"), 10, 64)
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid historyID",
			})
			return
		}

		gClient.SetEmail(email)
		gmailClient, err := gClient.GmailService(c)
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		historyListCall := gmailClient.Client.Users.History.List("me").StartHistoryId(historyID)

		historyList, err := historyListCall.Do()
		if err != nil {
			models.NewResponseInternalServerError(c)

			return
		}

		if len(historyList.History) == 0 {
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

		c.JSON(200, gin.H{
			"historyID": historyID,
			"messages":  messages,
		})
	}
}
