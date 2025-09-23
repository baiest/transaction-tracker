package main

import (
	"context"
	"log"
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"
	messageRepostiroy "transaction-tracker/internal/messages/repository"
	messageUsecase "transaction-tracker/internal/messages/usecase"
	movementRepostiroy "transaction-tracker/internal/movements/repository"
	movementUsecase "transaction-tracker/internal/movements/usecase"
	"transaction-tracker/logger"
	"transaction-tracker/pkg/databases/mongo"

	"transaction-tracker/pkg/databases/postgres"

	_ "transaction-tracker/env"
)

func main() {
	ctx := context.Background()

	dbClient, err := postgres.NewClient(ctx)
	if err != nil {
		log.Fatal("Unable to create postgres database client:", err)
	}

	defer dbClient.Close()

	client, err := mongo.NewClient(ctx)
	if err != nil {
		log.Fatal("Unable to create mongo database client:", err)
	}

	messageCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.MESSAGES)
	if err != nil {
		log.Fatal("Unable to get message collection:", err)
	}

	movementRepo := movementRepostiroy.NewPostgresRepository(dbClient.GetPool())
	movementUsecase := movementUsecase.NewMovementUsecase(movementRepo)
	movementHandler := handler.NewMovementHandler(movementUsecase)

	messageRepo := messageRepostiroy.NewMessageRepository(ctx, messageCollection)

	logger, err := logger.GetLogger(ctx, "messages-usecase")
	if err != nil {
		log.Fatal("Unable to create logger message-usecase:", err)
	}

	messageUsecase := messageUsecase.NewMessageUsecase(ctx, logger, messageRepo, movementRepo)
	messageHandler := handler.NewMessageHandler(messageUsecase)

	s := models.NewServer(8080)

	routerHandler := &routes.RouteHandler{
		MovementHandler: movementHandler,
		MessageHandler:  messageHandler,
	}

	s.AddRoutes(routerHandler.Routes())

	s.Run()
}
