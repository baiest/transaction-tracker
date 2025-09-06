package main

import (
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"

	_ "transaction-tracker/env"
)

func main() {
	s := models.NewServer(8080)

	s.AddRoutes(routes.Routes())

	s.Run()
}
