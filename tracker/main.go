package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "transaction-tracker/env"
	"transaction-tracker/googleapi"
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

const (
	STORE_EMAIL_MAX_RETRIES = 5
)

var (
	baseURL       = os.Getenv("BASE_TRANSACTION_URL")
	urlStoreEmail = "/api/v1/gmail/emails/%d/save"
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
		time.Sleep(5 * time.Second)
		return storeEmail(message, maxRetries-1)
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("Error unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func main() {
	ctx := context.Background()

	if urlStoreEmail == "" {
		log.Fatal("missing url to store the email")
	}

	pubsubService, err := googleapi.NewGooglePubSub(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	err = pubsubService.Subscribe(ctx, subscription, func(ctx context.Context, msg []byte) error {
		message := &Message{}
		err := json.Unmarshal(msg, message)
		if err != nil {
			return err
		}

		log.Printf("event: message_received, email: %s, history_id: %d", message.EmailAdress, message.HistoryID)

		err = storeEmail(message, STORE_EMAIL_MAX_RETRIES)
		if err != nil {
			log.Print("Error storing email:", err)

			return nil
		}

		log.Printf("event: message_stored, email: %s, history_id: %d", message.EmailAdress, message.HistoryID)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
