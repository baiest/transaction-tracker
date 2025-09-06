package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "transaction-tracker/env"
	"transaction-tracker/googleapi"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"
	"transaction-tracker/shared"
)

const (
	projectID    = "transaction-tracker-2473"
	subscription = "gmail-notifications-sub"
	topic        = "gmail-notifications"
)

type Message struct {
	EmailAdress string `json:"emailAddress"`
	HistoryID   uint64 `json:"historyId"`
}

func (m *Message) LogProperties() map[string]string {
	return map[string]string{
		"email":      m.EmailAdress,
		"history_id": fmt.Sprintf("%d", m.HistoryID),
	}
}

const (
	STORE_EMAIL_MAX_RETRIES = 5
)

var (
	baseURL         = os.Getenv("BASE_TRANSACTION_URL")
	timeToRetry     = 5 * time.Second
	urlStoreEmail   = "/api/v1/gmail/emails/histories/%d/save"
	credentialsFile = "sa-key.json"
	log             *loggerModels.Logger
)

func storeEmail(message *Message, maxRetries int) error {
	formData := fmt.Sprintf("email=%s", message.EmailAdress)

	req, err := http.NewRequest("POST", baseURL+fmt.Sprintf(urlStoreEmail, message.HistoryID), bytes.NewBufferString(formData))
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := shared.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == 404 && maxRetries > 0 {
		time.Sleep(timeToRetry)

		return storeEmail(message, maxRetries-1)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func handleSubscription(ctx context.Context, msg []byte) error {
	message := &Message{}
	err := json.Unmarshal(msg, message)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "error_unmarshalling_message",
			Error: err,
		})

		return err
	}

	log.Info(loggerModels.LogProperties{
		Event: "message_received",
		AdditionalParams: []loggerModels.Properties{
			message,
		},
	})

	err = storeEmail(message, STORE_EMAIL_MAX_RETRIES)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "error_storing_email",
			Error: err,
			AdditionalParams: []loggerModels.Properties{
				message,
			},
		})

		return nil
	}

	log.Info(loggerModels.LogProperties{
		Event: "message_stored",
		AdditionalParams: []loggerModels.Properties{
			message,
		},
	})

	return nil
}

func main() {
	ctx := context.Background()

	var err error

	log, err = logger.GetLogger(ctx, "transaction-tracker")
	if err != nil {
		fmt.Printf("Error getting logger: %v\n", err)

		return
	}

	if urlStoreEmail == "" {
		log.Error(loggerModels.LogProperties{
			Event: "missing_url_store_email",
			AdditionalParams: []loggerModels.Properties{
				logger.MapToProperties(map[string]string{
					"base_url": baseURL,
				}),
			},
		})

		return
	}

	pubsubService, err := googleapi.NewGooglePubSub(ctx, projectID, credentialsFile)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_initialize_pubsub",
			Error: err,
		})

		return
	}

	sub, err := pubsubService.GetSubscription(ctx, subscription)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_get_subscription",
			Error: err,
		})

		return
	}

	log.Info(loggerModels.LogProperties{
		Event: "pubsub_subscribed",
		AdditionalParams: []loggerModels.Properties{
			logger.MapToProperties(map[string]string{
				"subscription": sub.String(),
			}),
		},
	})

	err = pubsubService.Subscribe(ctx, subscription, handleSubscription)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_subscribe_pubsub",
			Error: err,
		})
	}
}
