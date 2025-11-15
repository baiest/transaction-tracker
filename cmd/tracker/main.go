package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	_ "transaction-tracker/env"
	accountsRepository "transaction-tracker/internal/accounts/repository"
	accountsUsecase "transaction-tracker/internal/accounts/usecase"
	extractsRepository "transaction-tracker/internal/extracts/repository"
	extractsUsecase "transaction-tracker/internal/extracts/usecase"
	messagesRepository "transaction-tracker/internal/messages/repository"
	messagesUsecase "transaction-tracker/internal/messages/usecase"
	movementsRepository "transaction-tracker/internal/movements/repository"
	movementsUsecase "transaction-tracker/internal/movements/usecase"
	"transaction-tracker/logger"
	loggerModels "transaction-tracker/logger/models"
	"transaction-tracker/pkg/databases/mongo"
	"transaction-tracker/pkg/databases/postgres"
	"transaction-tracker/pkg/google"
)

type subscriptionUsecase struct {
	accountsUsecase accountsUsecase.AccountsUsecase
	messagesUsecase messagesUsecase.MessageUsecase
}

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
	credentialsFile = "sa-key.json"
	log             *loggerModels.Logger
)

func (s *subscriptionUsecase) handleSubscription(ctx context.Context, msg []byte) error {
	time.Sleep(2 * time.Second)

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

	account, err := s.accountsUsecase.GetAccountByEmail(ctx, message.EmailAdress)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "error_getting_account",
			Error: err,
			AdditionalParams: []loggerModels.Properties{
				message,
			},
		})

		return err
	}

	err = s.accountsUsecase.RefreshGoogleToken(ctx, account)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "error_refreshing_google_token",
			Error: err,
			AdditionalParams: []loggerModels.Properties{
				message,
				account,
			},
		})

		return err
	}

	messages, err := s.messagesUsecase.ProcessByNotification(ctx, account, message.HistoryID)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "error_processing_messages",
			Error: err,
			AdditionalParams: []loggerModels.Properties{
				message,
			},
		})

		return err
	}

	if len(messages) == 0 {
		return nil
	}

	for _, m := range messages {
		log.Info(loggerModels.LogProperties{
			Event: "message_processed",
			AdditionalParams: []loggerModels.Properties{
				m,
			},
		})
	}

	log.Info(loggerModels.LogProperties{
		Event: "message_stored",
		AdditionalParams: []loggerModels.Properties{
			message,
		},
	})

	return nil
}

func NewSubscriptionsecase(ctx context.Context) (*subscriptionUsecase, error) {
	log := ctx.Value("logger").(*loggerModels.Logger)

	dbClient, err := postgres.NewClient(ctx)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_create_postgres_client",
			Error: err,
		})

		return nil, err
	}

	googleClient, err := google.NewGoogleClient(ctx)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_initialize_google_client",
			Error: err,
		})

		return nil, err
	}

	ctx, client, err := mongo.NewClient(ctx)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_create_mongo_client",
			Error: err,
		})

		return nil, err
	}

	accountCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.ACCOUNTS)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_get_account_collection",
			Error: err,
		})

		return nil, err
	}

	messageCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.MESSAGES)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_get_message_collection",
			Error: err,
		})

		return nil, err
	}

	extractsCollection, err := client.Collection(mongo.TRANSACTIONS, mongo.EXTRACTS)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_get_extracts_collection",
			Error: err,
		})

		return nil, err
	}

	accRepo := accountsRepository.NewAccountsRepository(accountCollection)
	accUsecase := accountsUsecase.NewAccountsUseCase(googleClient, accRepo)

	movementsRepo := movementsRepository.NewPostgresRepository(dbClient.GetPool())
	mvmUsecase := movementsUsecase.NewMovementUsecase(ctx, movementsRepo)

	extractsRepo := extractsRepository.NewExtractsRepository(extractsCollection)
	extractUsecase := extractsUsecase.NewExtractsUsecase(googleClient, extractsRepo)

	messageRepo := messagesRepository.NewMessageRepository(messageCollection)
	messageUsecase := messagesUsecase.NewMessageUsecase(ctx, googleClient, messageRepo, mvmUsecase, extractUsecase)

	return &subscriptionUsecase{
		accountsUsecase: accUsecase,
		messagesUsecase: messageUsecase,
	}, nil
}

func main() {
	ctx := context.Background()

	var err error

	log, err = logger.GetLogger(ctx, "tracker")
	if err != nil {
		fmt.Printf("Error getting logger: %v\n", err)

		return
	}

	ctx = context.WithValue(ctx, "logger", log)

	pubsubService, err := google.NewGooglePubSub(ctx, projectID, credentialsFile)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_initialize_pubsub",
			Error: err,
		})

		return
	}

	log.Info(loggerModels.LogProperties{
		Event:            "pubsub_subscribed",
		AdditionalParams: []loggerModels.Properties{},
	})

	s, err := NewSubscriptionsecase(ctx)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_create_subscription_usecase",
			Error: err,
		})

		return
	}

	err = pubsubService.Subscribe(ctx, subscription, s.handleSubscription)
	if err != nil {
		log.Error(loggerModels.LogProperties{
			Event: "failed_to_subscribe_pubsub",
			Error: err,
		})
	}
}
