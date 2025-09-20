package main

import (
	"context"
	"log"
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"
	"transaction-tracker/internal/movements/repository"
	"transaction-tracker/internal/movements/usecase"
	"transaction-tracker/pkg/databases/postgres"

	_ "transaction-tracker/env"
)

func main() {
	ctx := context.Background()

	dbClient, err := postgres.NewClient(ctx)
	if err != nil {
		log.Fatal("Unable to create database client:", err)
	}

	defer dbClient.Close()

	movementRepo := repository.NewPostgresRepository(dbClient.GetPool())
	movementUsecase := usecase.NewMovementUsecase(movementRepo)
	movementHandler := handler.NewMovementHandler(movementUsecase)

	s := models.NewServer(8080)

	routerHandler := &routes.RouteHandler{
		MovementHandler: movementHandler,
	}

	s.AddRoutes(routerHandler.Routes())

	s.Run()
}
