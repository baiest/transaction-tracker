package main

import (
	"context"
	"log"
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"

	"transaction-tracker/googleapi"

	_ "transaction-tracker/env"
)

func main() {
	gClient, err := googleapi.NewGoogleClient(context.Background())
	if err != nil {
		log.Println("Error creating google client:", err)

		return
	}

	s := models.NewServer(8080)

	s.AddRoutes(routes.Routes(), gClient)

	s.Run()
}
