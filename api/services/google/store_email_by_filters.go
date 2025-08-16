package google

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"transaction-tracker/api/models"
	"transaction-tracker/api/services/transactions"
	"transaction-tracker/googleapi"

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
			Subject: "test",
		},
	}
)

func StoreEmailByFilters(gClient *googleapi.GoogleClient) gin.HandlerFunc {
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
			models.NewResponseInternalServerError(c)
			return
		}

		if len(historyList.History) == 0 {
			models.NewResponseNotFoud(c, models.Response{
				Message: "no emails found for given historyID",
			})

			return
		}

		var messages []map[string]interface{}
		messageIDs := map[string]map[string]interface{}{}

		for _, h := range historyList.History {
			for _, m := range h.Messages {
				if _, ok := messageIDs[m.Id]; ok {
					continue
				}

				msg, err := gmailClient.Client.Users.Messages.Get("me", m.Id).Format("full").Do()
				if err != nil {
					continue
				}

				if !isMessageFiltered(msg) {
					continue
				}

				tr, err := parseEmailMessageToTransactionRequest(msg)
				if err != nil {
					continue
				}

				err = transactions.Create(tr)
				if err != nil {
					continue
				}

				message := map[string]interface{}{
					"id":      msg.Id,
					"payload": msg.Payload,
				}

				messages = append(messages, message)

				messageIDs[msg.Id] = message
			}
		}

		models.NewResponseOK(c, models.Response{
			Message: "messages",
			Data: map[string]any{
				"history_id": historyID,
				"messages":   messages,
			},
		})
	}
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
	if len(msg.Payload.Parts) == 0 {
		return nil, fmt.Errorf("missing body")
	}

	body := msg.Payload.Parts[0].Body.Data

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
