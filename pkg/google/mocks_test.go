package google

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockGoogleClient(t *testing.T) {
	c := require.New(t)

	mockClient := new(MockGoogleClient)

	account := &GoogleAccount{}
	mockClient.On("SaveTokenAndInitServices", context.Background(), "test_code").Return(account, nil)

	result, err := mockClient.SaveTokenAndInitServices(context.Background(), "test_code")

	c.NoError(err)
	c.Equal(account, result)
	mockClient.AssertExpectations(t)
}

func TestMockGmailService(t *testing.T) {
	c := require.New(t)
	mockService := new(MockGmailService)

	mockService.On("CreateWatch", context.Background(), "test_topic").Return(uint64(123), int64(456), nil)

	historyID, expiration, err := mockService.CreateWatch(context.Background(), "test_topic")

	c.NoError(err)
	c.Equal(uint64(123), historyID)
	c.Equal(int64(456), expiration)
	mockService.AssertExpectations(t)
}

func TestMockPubSubService(t *testing.T) {
	c := require.New(t)
	mockPubSub := new(MockPubSubService)

	mockPubSub.On("Publish", context.Background(), "test_topic", []byte("test_data")).Return(nil)

	err := mockPubSub.Publish(context.Background(), "test_topic", []byte("test_data"))

	c.NoError(err)
	mockPubSub.AssertExpectations(t)
}
