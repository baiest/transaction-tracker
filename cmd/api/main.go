package main

import (
	"context"
	"log"
	"transaction-tracker/api/handler"
	"transaction-tracker/api/models"
	"transaction-tracker/api/routes"
	accountRepository "transaction-tracker/internal/accounts/repository"
	accountUsecase "transaction-tracker/internal/accounts/usecase"
	extractRepostory "transaction-tracker/internal/extracts/repository"
	extractUsecase "transaction-tracker/internal/extracts/usecase"
	messageRepository "transaction-tracker/internal/messages/repository"
	messageUsecase "transaction-tracker/internal/messages/usecase"
	movementRepostiroy "transaction-tracker/internal/movements/repository"
	movementUsecase "transaction-tracker/internal/movements/usecase"
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

	ctx, client, err := mongo.NewClient(ctx)
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

	extractCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.EXTRACTS)
	if err != nil {
		log.Fatal("Unable to get account collection:", err)
	}

	movementRepo := movementRepostiroy.NewPostgresRepository(dbClient.GetPool())
	movementUsecase := movementUsecase.NewMovementUsecase(ctx, movementRepo)
	movementHandler := handler.NewMovementHandler(movementUsecase)

	googleClient, err := google.NewGoogleClient(ctx)
	if err != nil {
		log.Fatal("Unable to create google client:", err)
	}

	extractRepo := extractRepostory.NewExtractsRepository(extractCollection)
	extractUsecase := extractUsecase.NewExtractsUsecase(googleClient, extractRepo)

	messageRepo := messageRepository.NewMessageRepository(messageCollection)
	messageUsecase := messageUsecase.NewMessageUsecase(ctx, googleClient, messageRepo, movementUsecase, extractUsecase)
	messageHandler := handler.NewMessageHandler(messageUsecase)

	extractHandler := handler.NewExtractsHandler(extractUsecase, messageUsecase)

	accountRepo := accountRepository.NewAccountsRepository(accountCollection)
	accountUsecase := accountUsecase.NewAccountsUseCase(googleClient, accountRepo)
	accountHandler := handler.NewAccountHandler(accountUsecase)

	s := models.NewServer(accountUsecase, 8080)

	routerHandler := &routes.RouteHandler{
		AccountHandler:  accountHandler,
		MessageHandler:  messageHandler,
		ExtractHandler:  extractHandler,
		MovementHandler: movementHandler,
	}

	s.AddRoutes(routerHandler.Routes())

	s.Run()
}
