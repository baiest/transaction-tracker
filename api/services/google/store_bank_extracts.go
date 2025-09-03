package google

import (
	"fmt"
	"os"
	"transaction-tracker/api/models"
	movementsServices "transaction-tracker/api/services/movements"
	"transaction-tracker/database/mongo/schemas"
	documentextractor "transaction-tracker/document-extractor"
	"transaction-tracker/googleapi"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/gmail/v1"
)

var (
	daviviendaExtractor = &documentextractor.DaviviendaExtract{
		Password: os.Getenv("EXTRACT_PDF_PASSWORD"),
	}
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

		emailExtracts, err := downloadAttachments(c, log, gmailClient)
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		err = storeMovementsFromExtracts(c, log, emailExtracts)
		if err != nil {
			models.NewResponseInvalidRequest(c, models.Response{
				Message: err.Error(),
			})

			return
		}

		models.NewResponseOK(c, models.Response{
			Message: "success",
			Data:    emailExtracts,
		})
	}
}

func downloadAttachments(c *gin.Context, log *loggerModels.Logger, gmailClient *googleapi.GmailService) ([]*schemas.GmailExtract, error) {
	extracts, err := gmailClient.GetExtractMessages(c, "davivienda")
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "get_extract_messages_failed",
			Error: err,
		})

		return nil, err
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

			log.Info(loggerModels.LogProperties{
				Event: "download_attachments_success",
				AdditionalParams: []loggerModels.Properties{
					logger.MapToProperties(map[string]string{
						"message_id": msg.Id,
						"file_path":  extract.FilePath,
					}),
				},
			})
		}(msg)
	}

	for range extracts.Messages {
		<-messages
	}

	close(messages)

	if len(downloadsFailed) > 0 {
		log.Error(loggerModels.LogProperties{
			Event: "download_attachments_failed",
			Error: fmt.Errorf("failed to download attachments: %v", downloadsFailed),
		})

		return nil, fmt.Errorf("failed to download attachments: %v", downloadsFailed)
	}

	return emailExtracts, nil
}

func storeMovementsFromExtracts(c *gin.Context, log *loggerModels.Logger, emailExtracts []*schemas.GmailExtract) error {
	movementsService, err := movementsServices.NewMovementsService(c)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "init_movements_service_failed",
			Error: err,
		})

		return err
	}

	movementsChan := make(chan bool, len(emailExtracts))

	for _, extract := range emailExtracts {
		go func(extract *schemas.GmailExtract) {
			defer func() { movementsChan <- true }()

			movements := daviviendaExtractor.GetMovements("api\\" + extract.FilePath)

			mChan := make(chan bool, len(movements))

			for _, m := range movements {
				go func(m *documentextractor.Movement) {
					defer func() { mChan <- true }()

					movement := schemas.NewMovement(
						extract.Email,
						extract.ID,
						m.Date,
						m.Value,
						m.IsNegative,
						m.Type,
						m.Detail,
					)

					err := movementsService.CreateMovement(c, movement)
					if err != nil {
						log.Error(loggerModels.LogProperties{
							Event: "create_movement_failed",
							Error: err,
						})

						return
					}

					log.Info(loggerModels.LogProperties{
						Event: "create_movement_success",
					})
				}(m)
			}

			for range movements {
				<-mChan
			}

			close(mChan)
		}(extract)
	}

	for range emailExtracts {
		<-movementsChan
	}

	close(movementsChan)

	return nil
}
