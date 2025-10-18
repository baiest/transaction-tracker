package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	gmailv1 "google.golang.org/api/gmail/v1"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	extractUsecase "transaction-tracker/internal/extracts/usecase"
	"transaction-tracker/internal/messages/domain"
	repo "transaction-tracker/internal/messages/repository"
	movementsUsecase "transaction-tracker/internal/movements/usecase"
	"transaction-tracker/pkg/google"
)

func TestProcess_MissingGmailService(t *testing.T) {
	c := require.New(t)

	mockGoogle := &google.MockGoogleClient{}
	mockGoogle.On("GmailService", mock.Anything, mock.Anything).Return(nil, errors.New("no gmail"))

	mockMessages := new(repo.MockMessageRepository)
	mockMessages.On("GetMessageByExternalID", mock.Anything, "ext-1", "acc1").Return(nil, nil)

	mockMovementsUsecase := new(movementsUsecase.MockMovementUsecase)

	mockExtractsUsecase := new(extractUsecase.MockExtractsUsecase)
	mockExtractsUsecase.On("GetByMessageID", mock.Anything, "msg-1").Return(nil, nil)
	mockExtractsUsecase.On("Update", mock.Anything, mock.Anything).Return(nil)

	u := NewMessageUsecase(context.Background(), mockGoogle, mockMessages, mockMovementsUsecase, mockExtractsUsecase)

	account := &accountsDomain.Account{ID: "acc1", GoogleAccount: &google.GoogleAccount{}}
	msg, err := u.Process(context.Background(), "notif1", "ext-1", account)
	c.Nil(msg)
	c.Error(err)
	mockGoogle.AssertCalled(t, "GmailService", mock.Anything, mock.Anything)
}

func TestProcess_MessageAlreadyExists_ReturnsExisting(t *testing.T) {
	c := require.New(t)

	// Gmail message returned by service
	gmsg := &gmailv1.Message{
		Id: "ext-1",
		Payload: &gmailv1.MessagePart{
			Headers: []*gmailv1.MessagePartHeader{
				{Name: "From", Value: "bancodavivienda@davivienda.com"},
				{Name: "Subject", Value: "Extractos Septiembre"},
			},
		},
	}

	// mock gmail service
	mockGmail := &google.MockGmailService{}
	mockGmail.On("GetMessageByID", mock.Anything, "ext-1").Return(gmsg, nil)

	// mock google client that returns the mock gmail service
	mockGoogle := &google.MockGoogleClient{}
	mockGoogle.On("GmailService", mock.Anything, mock.Anything).Return(mockGmail, nil)

	// repository mock (testify/mock) that simulates an existing DB message
	repoMock := &repo.MockMessageRepository{}
	existing := &domain.Message{ID: "msg-db-1", ExternalID: "ext-1", AccountID: "acc1"}
	repoMock.On("GetMessageByExternalID", mock.Anything, "ext-1", "acc1").Return(existing, nil)

	u := &messageUsecase{
		messageRepo:    repoMock,
		mvmUsecase:     nil,
		extractUsecase: new(extractUsecase.MockExtractsUsecase),
		googleClient:   mockGoogle,
		log:            nil,
	}

	account := &accountsDomain.Account{ID: "acc1", GoogleAccount: &google.GoogleAccount{}}
	got, err := u.Process(context.Background(), "notif1", "ext-1", account)
	c.NoError(err)
	c.Equal(existing, got)

	repoMock.AssertCalled(t, "GetMessageByExternalID", mock.Anything, "ext-1", "acc1")
	repoMock.AssertNotCalled(t, "SaveMessage", mock.Anything, mock.Anything)
	mockGoogle.AssertCalled(t, "GmailService", mock.Anything, mock.Anything)
	mockGmail.AssertCalled(t, "GetMessageByID", mock.Anything, "ext-1")
}

func TestProcess_NewMessage_SaveCalled(t *testing.T) {
	c := require.New(t)

	gmsg := &gmailv1.Message{
		Id: "ext-2",
		Payload: &gmailv1.MessagePart{
			Headers: []*gmailv1.MessagePartHeader{
				{Name: "From", Value: "banco_davivienda@davivienda.com"},
				{Name: "Subject", Value: "Notificación de movimiento"},
			},
		},
	}

	mockGmail := &google.MockGmailService{}
	mockGmail.On("GetMessageByID", mock.Anything, "ext-2").Return(gmsg, nil)

	mockGoogle := &google.MockGoogleClient{}
	mockGoogle.On("GmailService", mock.Anything, mock.Anything).Return(mockGmail, nil)

	repoMock := &repo.MockMessageRepository{}
	repoMock.On("GetMessageByExternalID", mock.Anything, "ext-2", "acc2").Return(nil, repo.ErrMessageNotFound)
	repoMock.On("SaveMessage", mock.Anything, mock.AnythingOfType("*domain.Message")).Return(nil)

	u := &messageUsecase{
		messageRepo:    repoMock,
		mvmUsecase:     nil,
		extractUsecase: new(extractUsecase.MockExtractsUsecase),
		googleClient:   mockGoogle,
		log:            nil,
	}

	account := &accountsDomain.Account{ID: "acc2", GoogleAccount: &google.GoogleAccount{}}
	got, err := u.Process(context.Background(), "notif2", "ext-2", account)
	c.NoError(err)
	c.NotNil(got)

	repoMock.AssertCalled(t, "GetMessageByExternalID", mock.Anything, "ext-2", "acc2")
	repoMock.AssertCalled(t, "SaveMessage", mock.Anything, mock.AnythingOfType("*domain.Message"))
	mockGoogle.AssertCalled(t, "GmailService", mock.Anything, mock.Anything)
	mockGmail.AssertCalled(t, "GetMessageByID", mock.Anything, "ext-2")
}
