package main

import (
	"context"
	"log"
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"
	accountRepository "transaction-tracker/internal/accounts/repository"
	accountUsecase "transaction-tracker/internal/accounts/usecase"
	messageRepository "transaction-tracker/internal/messages/repository"
	messageUsecase "transaction-tracker/internal/messages/usecase"
	movementRepostiroy "transaction-tracker/internal/movements/repository"
	movementUsecase "transaction-tracker/internal/movements/usecase"
	"transaction-tracker/logger"
	"transaction-tracker/pkg/databases/mongo"
	"transaction-tracker/pkg/google"

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

	accountCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.ACCOUNTS)
	if err != nil {
		log.Fatal("Unable to get account collection:", err)
	}

	movementRepo := movementRepostiroy.NewPostgresRepository(dbClient.GetPool())
	movementUsecase := movementUsecase.NewMovementUsecase(movementRepo)
	movementHandler := handler.NewMovementHandler(movementUsecase)

	messageRepo := messageRepository.NewMessageRepository(messageCollection)

	logger, err := logger.GetLogger(ctx, "messages-usecase")
	if err != nil {
		log.Fatal("Unable to create logger message-usecase:", err)
	}

	googleClient, err := google.NewGoogleClient(ctx)
	if err != nil {
		log.Fatal("Unable to create google client:", err)
	}

	gmailClient, err := google.NewGmailClient(ctx, googleClient)
	if err != nil {
		log.Fatal("Unable to create gmail client:", err)
	}

	messageUsecase := messageUsecase.NewMessageUsecase(ctx, logger, googleClient, messageRepo, movementRepo)
	messageHandler := handler.NewMessageHandler(messageUsecase)

	accountRepo := accountRepository.NewAccountsRepository(accountCollection)
	accountUsecase := accountUsecase.NewAccountsUseCase(googleClient, gmailClient, accountRepo)
	accountHandler := handler.NewAccountHandler(accountUsecase)

	s := models.NewServer(accountUsecase, 8080)

	routerHandler := &routes.RouteHandler{
		AccountHandler:  accountHandler,
		MovementHandler: movementHandler,
		MessageHandler:  messageHandler,
	}

	s.AddRoutes(routerHandler.Routes())

	s.Run()
}
