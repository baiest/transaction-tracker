package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MessageStatus string

const (
	_message_prefix = "MSI"

	Pending MessageStatus = "pending"
	Success MessageStatus = "success"
	Failure MessageStatus = "failure"
)

type Message struct {
	ID             string        `bson:"_id" json:"id"`
	AccountID      string        `bson:"account_id" json:"account_id"`
	ExternalID     string        `bson:"external_id" json:"external_id"`
	NotificationID string        `bson:"notification_id,omitempty" json:"notification_id,omitempty"`
	ExtractID      string        `bson:"extract_id,omitempty" json:"extract_id,omitempty"`
	From           string        `bson:"from,omitempty" json:"from,omitempty"`
	To             string        `bson:"to,omitempty" json:"to,omitempty"`
	Status         MessageStatus `bson:"status" json:"status"`
	FailureReason  string        `bson:"failure_reason" json:"failure_reason"`
	Date           time.Time     `bson:"date" json:"date"`
	CreatedAt      time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at" json:"updated_at"`
}

// NewMessage creates a new Message with sensible defaults
func NewMessage(accountID, from string, to string, externalID string, notificationID string, extractID string, date time.Time) *Message {
	return &Message{
		ID:             _message_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		AccountID:      accountID,
		ExternalID:     externalID,
		NotificationID: notificationID,
		ExtractID:      extractID,
		From:           from,
		To:             to,
		Status:         Pending,
		Date:           date,
	}
}
