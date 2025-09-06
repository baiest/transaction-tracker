package google

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"transaction-tracker/api/models"
	movementsServices "transaction-tracker/api/services/movements"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi"
	"transaction-tracker/googleapi/repositories"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/gmail/v1"
)

type EmailFilter struct {
	Email   string
	Subject string
}

var (
	emailFilters = []EmailFilter{
		{
			Email:   "banco_davivienda@davivienda.com",
			Subject: "davivienda",
		},
		{
			Email:   "juanballesteros2001@gmail.com",
			Subject: "davivienda",
		},
	}

	log *loggerModels.Logger
)

func StoreEmailByFilters() gin.HandlerFunc {
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

		historyListCall := gmailClient.Client.Users.History.List("me").StartHistoryId(historyID)

		historyList, err := historyListCall.Do()
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "get_history_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		if len(historyList.History) == 0 {
			log.Error(loggerModels.LogProperties{
				Event: "history_id_not_found",
				Error: err,
			})

			models.NewResponseNotFound(c, models.Response{
				Message: "no emails found for given historyID",
			})

			return
		}

		notification, err := gmailClient.SaveNotification(c, strconv.FormatInt(int64(historyID), 10))
		if err != err {
			log.Error(loggerModels.LogProperties{
				Event: "save_notification_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		if notification.Status != "pending" {
			processNotificationMessages(c, log, gmailClient, notification)

			return
		}

		var messages []map[string]interface{}
		messageIDs := map[string]bool{}

		notificationStatus := "pending"

		for _, h := range historyList.History {
			stop := make(chan bool, len(h.Messages))

			for _, m := range h.Messages {
				if messageIDs[m.Id] {
					stop <- true
					continue
				}

				messageIDs[m.Id] = true

				message, err := gmailClient.GetMessage(c, m.Id)
				if err != nil && !errors.Is(err, repositories.ErrMessageNotFound) {
					log.Error(loggerModels.LogProperties{
						Event: "get_message_by_id_failed",
						Error: err,
					})

					models.NewResponseInternalServerError(c)

					return
				}

				if message != nil {
					stop <- true
					continue
				}

				notificationMessage := &schemas.Message{ID: m.Id, Status: "pending", NotificationID: notification.ID}

				go func(messageID string, message *schemas.Message) {
					defer func() {
						stop <- true
					}()

					msg, err := processMessageInTransactions(c, gmailClient, messageID)
					if err != nil {
						log.Error(loggerModels.LogProperties{
							Event: "process_message_failed",
							Error: err,
						})

						notificationStatus = "failure"
						message.Status = "failure"

						err := gmailClient.SaveMessage(c, message)
						if err != nil {
							log.Error(loggerModels.LogProperties{
								Event: "save_message_status_failure_failed",
								Error: err,
							})
						}

						return
					}

					if msg == nil {
						log.Info(loggerModels.LogProperties{
							Event: "message_filtered",
							AdditionalParams: []loggerModels.Properties{
								logger.MapToProperties(map[string]string{
									"message_id": messageID,
								}),
							},
						})

						return
					}

					log.Info(loggerModels.LogProperties{
						Event: "movement_created",
						AdditionalParams: []loggerModels.Properties{
							logger.MapToProperties(map[string]string{
								"message_id": messageID,
							}),
						},
					})

					notificationMessage.Status = "success"
					err = gmailClient.SaveMessage(c, notificationMessage)
					if err != nil {
						log.Error(loggerModels.LogProperties{
							Event: "process_message_status_success_failed",
							Error: err,
						})

						notificationStatus = "failure"
					}

					messageResponse := map[string]interface{}{
						"id":      msg.Id,
						"payload": msg.Payload,
					}

					messages = append(messages, messageResponse)
				}(m.Id, notificationMessage)
			}

			for range h.Messages {
				<-stop
			}

			close(stop)
		}

		if notificationStatus == "pending" {
			notificationStatus = "success"
		}

		notification.Status = notificationStatus

		err = gmailClient.UpdateNotification(c, notification)
		if err != nil {
			log.Error(loggerModels.LogProperties{
				Event: "processupdate_notification_failed",
				Error: err,
			})

			models.NewResponseInternalServerError(c)

			return
		}

		models.NewResponseOK(c, models.Response{
			Data: map[string]any{
				"history_id": historyID,
				"messages":   messages,
			},
		})
	}
}

func processMessageInTransactions(c *gin.Context, gmailClient *googleapi.GmailService, messageID string) (*gmail.Message, error) {
	msg, err := gmailClient.Client.Users.Messages.Get("me", messageID).Format("full").Do()
	if err != nil {
		return nil, nil
	}

	if !isMessageFiltered(msg) {
		return nil, nil
	}

	tr, err := parseEmailMessageToTransactionRequest(msg)
	if err != nil {
		return nil, err
	}

	account := c.MustGet("account").(*models.Account)

	movement := schemas.NewMovement(
		account.Email,
		msg.Id,
		tr.Date,
		float64(tr.Value),
		float64(tr.Value) < 0,
		tr.Type,
		"",
	)

	movementsService, err := movementsServices.NewMovementsService(c)
	if err != nil {
		return nil, err
	}

	err = movementsService.CreateMovement(c, movement)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func processNotificationMessages(c *gin.Context, log *loggerModels.Logger, gmailClient *googleapi.GmailService, notification *schemas.GmailNotification) {
	if notification.Status == "success" {
		models.NewResponseOK(c, models.Response{
			Message: "notification is already saved",
		})

		return
	}

	if len(notification.Messages) == 0 {
		log.Info(loggerModels.LogProperties{
			Event: "notifications_empty",
		})

		models.NewResponseOK(c, models.Response{
			Message: "no messages found for notification",
		})

		return
	}

	messageByID := map[string]bool{}
	notificationStatus := "pending"

	stop := make(chan bool, len(notification.Messages))

	for _, message := range notification.Messages {
		if messageByID[message.ID] {
			continue
		}

		go func() {
			defer func() {
				stop <- true
			}()

			if message.Status == "success" {
				return
			}

			messageByID[message.ID] = true

			_, err := processMessageInTransactions(c, gmailClient, message.ID)
			if err != nil {
				log.Error(loggerModels.LogProperties{
					Event: "process_message_failed",
					Error: err,
				})

				notificationStatus = "failure"
				message.Status = "failure"

				updateErr := gmailClient.UpdateMessage(c, message)
				if updateErr != nil {
					log.Error(loggerModels.LogProperties{
						Event: "update_message_status_failure_failed",
						Error: err,
					})

					notificationStatus = "failure"
				}

				return
			}

			message.Status = "success"

			updateErr := gmailClient.UpdateMessage(c, message)
			if updateErr != nil {
				log.Error(loggerModels.LogProperties{
					Event: "update_message_status_succes_failed",
					Error: updateErr,
				})

				notificationStatus = "failure"
			}
		}()
	}

	for range notification.Messages {
		<-stop
	}

	close(stop)

	if notificationStatus == "pending" {
		notificationStatus = "success"
	}

	notification.Status = notificationStatus

	updateNotificationErr := gmailClient.UpdateNotification(c, notification)
	if updateNotificationErr != nil {
		log.Error(loggerModels.LogProperties{
			Event: "update_notification_failed",
			Error: updateNotificationErr,
		})

		models.NewResponseInternalServerError(c)

		return
	}

	models.NewResponseOK(c, models.Response{
		Data: map[string]any{
			"history_id": notification.ID,
			"messages":   notification.Messages,
		},
	})
}

// Review this
func isMessageFiltered(msg *gmail.Message) bool {
	var from, subject string
	if msg.Payload != nil && msg.Payload.Headers != nil {
		for _, header := range msg.Payload.Headers {
			if header.Name == "From" {
				from = header.Value
			}

			if header.Name == "Subject" {
				subject = header.Value
			}
		}
	}

	for _, filter := range emailFilters {
		if strings.Contains(strings.ToLower(from), filter.Email) && strings.Contains(strings.ToLower(subject), filter.Subject) {
			return true
		}
	}

	return false
}

func parseEmailMessageToTransactionRequest(msg *gmail.Message) (*models.TransactionRequest, error) {
	body := ""

	if len(msg.Payload.Parts) > 0 {
		body = msg.Payload.Parts[0].Body.Data
	} else {
		body = msg.Payload.Body.Data
	}

	if body == "" {
		return nil, fmt.Errorf("missing body")
	}

	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return nil, fmt.Errorf("Error decoding body")
	}

	dm, err := models.NewDaviviendaMovementFromText(string(decodedBody))
	if err != nil {
		return nil, err
	}

	return dm.ToTransactionRequest(), nil
}
