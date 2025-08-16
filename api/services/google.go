package services

import (
	"errors"
	"fmt"
	"strconv"
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

func GoogleLogin(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := gClient.SaveTokenAndInitServices(c, c.Query("code"))
		if err != nil {
			fmt.Println("Error saving token:", err)

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
			fmt.Println("Error creating watch:", err)

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: fmt.Sprintf("HistoryID: %d ExpirationTime: %d", historyID, expirationTime),
		})
	}
}

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
			fmt.Println("Error deleting watcher:", err)

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "watch deleted succefully",
		})
	}
}

func GetEmailByHistoryID(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		historyID, err := strconv.ParseUint(c.Param("historyID"), 10, 64)
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid historyID",
			})
			return
		}

		gmailClient, err := gClient.GmailService()
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		historyListCall := gmailClient.Client.Users.History.List("me").StartHistoryId(historyID)

		historyList, err := historyListCall.Do()
		if err != nil {
			fmt.Println("Error fetching history:", err)
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
					fmt.Println("Error fetching message:", err)
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
