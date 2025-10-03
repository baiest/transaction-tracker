package google

import (
	"context"
	"time"

	pubsub "cloud.google.com/go/pubsub"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

// MockGoogleClient is a mock of GoogleClientAPI
type MockGoogleClient struct {
	mock.Mock
}

func (m *MockGoogleClient) SaveTokenAndInitServices(ctx context.Context, code string) (*GoogleAccount, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*GoogleAccount), args.Error(1)
}

func (m *MockGoogleClient) GetAuthURL() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockGoogleClient) SetToken(token *oauth2.Token) {
	m.Called(token)
}

func (m *MockGoogleClient) GetUserEmail(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockGoogleClient) GmailService(ctx context.Context, googleAccount *GoogleAccount) (GmailAPI, error) {
	args := m.Called(ctx, googleAccount)
	return args.Get(0).(GmailAPI), args.Error(1)
}

func (m *MockGoogleClient) RefreshToken(ctx context.Context, googleAccount *GoogleAccount) (*oauth2.Token, error) {
	args := m.Called(ctx, googleAccount)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m *MockGoogleClient) Config() *oauth2.Config {
	args := m.Called()
	return args.Get(0).(*oauth2.Config)
}

// MockGmailService is a mock of GmailAPI
type MockGmailService struct {
	mock.Mock
}

func (m *MockGmailService) CreateWatch(ctx context.Context, topicName string) (uint64, int64, error) {
	args := m.Called(ctx, topicName)
	return args.Get(0).(uint64), args.Get(1).(int64), args.Error(2)
}

func (m *MockGmailService) DeleteWatch() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockGmailService) GetMessageByID(ctx context.Context, messageID string) (*gmail.Message, error) {
	args := m.Called(ctx, messageID)
	return args.Get(0).(*gmail.Message), args.Error(1)
}

func (m *MockGmailService) GetMessageAttachment(ctx context.Context, messageID string, attachmentID string) (*gmail.MessagePartBody, error) {
	args := m.Called(ctx, messageID, attachmentID)
	return args.Get(0).(*gmail.MessagePartBody), args.Error(1)
}

func (m *MockGmailService) GetExtractMessages(ctx context.Context, bankName string) (*gmail.ListMessagesResponse, error) {
	args := m.Called(ctx, bankName)
	return args.Get(0).(*gmail.ListMessagesResponse), args.Error(1)
}

func (m *MockGmailService) DownloadAttachments(ctx context.Context, accountID string, messageID string) (time.Month, int, string, error) {
	args := m.Called(ctx, accountID, messageID)
	return args.Get(0).(time.Month), args.Get(1).(int), args.String(2), args.Error(3)
}

// MockPubSubService is a mock of PubSubService
type MockPubSubService struct {
	mock.Mock
}

func (m *MockPubSubService) Publish(ctx context.Context, topic string, data []byte) error {
	args := m.Called(ctx, topic, data)
	return args.Error(0)
}

func (m *MockPubSubService) Subscribe(ctx context.Context, subscription string, handler func(ctx context.Context, msg []byte) error) error {
	args := m.Called(ctx, subscription, handler)
	return args.Error(0)
}

func (m *MockPubSubService) GetSubscription(ctx context.Context, subscriptionID string) (*pubsub.Subscription, error) {
	args := m.Called(ctx, subscriptionID)
	return args.Get(0).(*pubsub.Subscription), args.Error(1)
}
