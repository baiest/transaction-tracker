package usecase

import (
	"context"
	"errors"
	"testing"
	"time"
	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/extracts/domain"
	"transaction-tracker/internal/extracts/repository"
	"transaction-tracker/pkg/google"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/gmail/v1"
)

func TestGetByMessageID_Success(t *testing.T) {
	c := require.New(t)

	ex := domain.NewExtract("acc-1", "msg-1", "inst-1", "/p.pdf", time.March, 2025)

	mockRepo := new(repository.MockExtractsRepository)

	mockRepo.On("GetByMessageID", context.Background(), ex.MessageID).Return(ex, nil)

	googleClientMock := new(google.MockGoogleClient)

	u := NewExtractsUsecase(googleClientMock, mockRepo)

	res, err := u.GetByMessageID(context.Background(), ex.MessageID)
	c.NoError(err)
	c.Equal(ex, res)
}

func TestGetByMessageID_NotFound(t *testing.T) {
	c := require.New(t)

	mockRepo := new(repository.MockExtractsRepository)

	mockRepo.On("GetByMessageID", context.Background(), "noexist").Return(nil, repository.ErrExtractNotFound)

	googleClientMock := new(google.MockGoogleClient)

	u := NewExtractsUsecase(googleClientMock, mockRepo)

	res, err := u.GetByMessageID(context.Background(), "noexist")
	c.Nil(res)
	c.Error(err)
	c.Equal(ErrExtractNotFound, err)
}

func TestGetExtractMessages_Success(t *testing.T) {
	c := require.New(t)

	msgs := &gmail.ListMessagesResponse{
		Messages: []*gmail.Message{
			{Id: "m1"},
			{Id: "m2"},
		},
	}

	account := &accountsDomain.Account{
		ID:            "acc-1",
		GoogleAccount: &google.GoogleAccount{},
	}

	mockGmail := new(google.MockGmailService)
	mockGmail.On("GetExtractMessages", context.Background(), "SomeBank").Return(msgs, nil)

	mockGoogle := new(google.MockGoogleClient)
	mockGoogle.On("GmailService", context.Background(), account.GoogleAccount).Return(mockGmail, nil)

	mockRepo := new(repository.MockExtractsRepository)
	mockRepo.On("Save", context.Background(), &domain.Extract{}).Return(nil)

	u := NewExtractsUsecase(mockGoogle, mockRepo)

	res, err := u.GetExtractMessages(context.Background(), "SomeBank", account)
	c.NoError(err)
	c.Equal([]string{"m1", "m2"}, res)
}

func TestGetExtractMessages_GmailServiceError(t *testing.T) {
	c := require.New(t)

	mockGoogle := new(google.MockGoogleClient)
	mockGoogle.On("GmailService", context.Background(), &google.GoogleAccount{}).Return(nil, errors.New("gmail service failure"))

	account := &accountsDomain.Account{
		ID:            "acc-1",
		GoogleAccount: &google.GoogleAccount{},
	}

	mockRepo := new(repository.MockExtractsRepository)

	u := NewExtractsUsecase(mockGoogle, mockRepo)

	_, err := u.GetExtractMessages(context.Background(), "Bank", account)
	c.Error(err)
}

func TestGetExtractMessages_GetMessagesError(t *testing.T) {
	c := require.New(t)

	mockGmail := new(google.MockGmailService)
	mockGmail.On("GetExtractMessages", context.Background(), "Bank").Return(nil, errors.New("get messages failed"))

	mockGoogle := new(google.MockGoogleClient)
	mockGoogle.On("GmailService", mock.Anything, mock.Anything).Return(mockGmail, nil)

	account := &accountsDomain.Account{
		ID:            "acc-1",
		GoogleAccount: &google.GoogleAccount{},
	}

	mockRepo := new(repository.MockExtractsRepository)

	u := NewExtractsUsecase(mockGoogle, mockRepo)

	_, err := u.GetExtractMessages(context.Background(), "Bank", account)
	c.Error(err)
}

func TestSave_DelegatesToRepo(t *testing.T) {
	c := require.New(t)

	mockRepo := new(repository.MockExtractsRepository)

	mockRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

	mockGoogleClient := new(google.MockGoogleClient)

	u := NewExtractsUsecase(mockGoogleClient, mockRepo)

	err := u.Save(context.Background(), &domain.Extract{})
	c.NoError(err)

	called := mockRepo.AssertCalled(t, "Save", context.Background(), &domain.Extract{})

	c.True(called)
}

func TestUpdate_DelegatesToRepo(t *testing.T) {
	c := require.New(t)

	mockRepo := new(repository.MockExtractsRepository)
	mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil)

	mockGmail := new(google.MockGmailService)
	mockGmail.On("GetExtractMessages", context.Background(), "Bank").Return(&gmail.ListMessagesResponse{}, nil)

	mockGoogleClient := new(google.MockGoogleClient)
	mockGoogleClient.On("GmailService").Return(mockGmail, nil)

	u := NewExtractsUsecase(mockGoogleClient, mockRepo)

	err := u.Update(context.Background(), &domain.Extract{})
	c.NoError(err)

	called := mockRepo.AssertCalled(t, "Update", context.Background(), &domain.Extract{})

	c.True(called)
}
