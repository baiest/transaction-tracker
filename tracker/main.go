package main

import (
	"context"
	"fmt"
	"log"

	"transaction-tracker/googleapi"
)

const (
	projectID    = "transaction-tracker-2473"
	subscription = "gmail-notifications-sub"
	topic        = "gmail-notifications"
)

func main() {
	ctx := context.Background()

	pubsubService, err := googleapi.NewGooglePubSub(ctx, projectID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Esperando mensajes...")

	err = pubsubService.Subscribe(ctx, subscription, func(ctx context.Context, msg []byte) error {
		fmt.Println("Mensaje recibido:", string(msg))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
