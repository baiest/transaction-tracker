package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	c := require.New(t)

	accountID := "acc-123"
	from := "test@example.com"
	to := "dest@example.com"
	externalID := "ext-123"
	notificationID := "notif-456"
	extractID := "extr-789"

	now := time.Now()
	msg := NewMessage(
		accountID,
		from,
		to,
		externalID,
		notificationID,
		extractID,
		now,
	)

	c.Equal(accountID, msg.AccountID)
	c.Equal(from, msg.From)
	c.Equal(to, msg.To)
	c.Equal(Pending, msg.Status)

	c.Equal(externalID, msg.ExternalID)
	c.Equal(notificationID, msg.NotificationID)
	c.Equal(extractID, msg.ExtractID)

	c.WithinDuration(now, msg.Date, time.Second)
}
