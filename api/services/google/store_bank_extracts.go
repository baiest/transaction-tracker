package google

import (
	"fmt"
	"transaction-tracker/api/models"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/gmail/v1"
)

func StoreBankExtracts(gClient *googleapi.GoogleClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		log, err := logger.GetLogger(c, "transaction-tracker")
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: fmt.Sprintf("logger not init: %s", err.Error()),
			})

			return
		}

		email := c.PostForm("email")
		if email == "" {
			log.Info(loggerModels.LogProperties{
				Event: "missing_email",
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "email is required in x-www-form-urlencoded body",
			})

			return
		}

		gClient.SetEmail(email)

		gmailClient, err := gClient.GmailService(c)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "init_gmail_client_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})
			return
		}

		extracts, err := gmailClient.GetExtractMessages(c, "davivienda")
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_extract_messages_failed",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})
			return
		}

		downloadsFailed := []string{}
		messages := make(chan bool, len(extracts.Messages))
		emailExtracts := []*schemas.GmailExtract{}

		for _, msg := range extracts.Messages {
			go func(msg *gmail.Message) {
				defer func() { messages <- true }()

				extract, err := gmailClient.DownloadAttachments(c, msg.Id)
				if err != nil {
					log.Error(loggerModels.LogProperties{
						Event: "download_attachments_failed",
						Error: err,
					})

					downloadsFailed = append(downloadsFailed, msg.Id)

					return
				}

				emailExtracts = append(emailExtracts, extract)
			}(msg)
		}

		for i := 0; i < len(extracts.Messages); i++ {
			<-messages
		}

		if len(downloadsFailed) > 0 {
			log.Error(loggerModels.LogProperties{
				Event: "download_attachments_failed",
				Error: fmt.Errorf("failed to download attachments: %v", downloadsFailed),
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: fmt.Sprintf("failed to download attachments: %v", downloadsFailed),
			})

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "success",
			Data:    emailExtracts,
		})
	}
}
