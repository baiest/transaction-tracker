package gmail

import (
	"errors"
	"strconv"
	"transaction-tracker/api/models"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
)

func StoreEmailByFilters() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := c.MustGet("logger").(*loggerModels.Logger)

		historyID, err := strconv.ParseUint(c.Param("historyID"), 10, 64)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "invalid_history_id",
				Error: err,
			})

			models.NewResponseInvalidRequest(c, models.Response{
				Message: "invalid historyID",
			})

			return
		}

		gmailService, err := NewGmailService(c)
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

		notification, err := gmailService.CreateNotification(c, historyID)
		if errors.Is(err, ErrHistoryNotFound) {
			models.NewResponseNotFound(c, models.Response{
				Message: "history not found",
			})

			return
		}

		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "create_notification_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		if notification.Status == "success" {
			models.NewResponseOK(c, models.Response{
				Data: notification,
			})

			return
		}

		strHistoiryID := strconv.FormatInt(int64(historyID), 10)

		messages := []*schemas.Message{}

		notification.Status = "pending"

		messagesChan := make(chan bool, len(notification.Messages))

		for _, notificationMessage := range notification.Messages {
			go func() {
				defer func() {
					messagesChan <- true
				}()

				if notificationMessage.Status == "success" {
					return
				}

				message, err := gmailService.ProcessMessage(c, notificationMessage.ID, strHistoiryID)
				if err != nil {
					notificationMessage.Status = "failure"

					log.Error(loggerModels.LogProperties{
						Event: "process_message_failed",
						Error: err,
						AdditionalParams: []loggerModels.Properties{
							logger.MapToProperties(map[string]string{
								"history_id": strHistoiryID,
							}),
						},
					})

					if notification.Status == "pending" {
						notification.Status = "failure"

						err = gmailService.UpdateNotification(c, notification)
						if err != nil {
							log.Error(loggerModels.LogProperties{
								Event: "processupdate_notification_failed",
								Error: err,
							})
						}
					}

					if notificationMessage != nil {
						messages = append(messages, notificationMessage)
					}

					return
				}

				if message != nil && message.Status != "" {
					messages = append(messages, message)
				}
			}()
		}

		for range notification.Messages {
			<-messagesChan
		}

		close(messagesChan)

		notification.Messages = messages

		if notification.Status == "pending" {
			notification.Status = "success"

			err = gmailService.UpdateNotification(c, notification)
			if err != nil {
				log.Error(loggerModels.LogProperties{
					Event: "processupdate_notification_failed",
					Error: err,
				})

				models.NewResponseInternalServerError(c)

				return
			}
		}

		models.NewResponseOK(c, models.Response{
			Data: notification,
		})
	}
}
