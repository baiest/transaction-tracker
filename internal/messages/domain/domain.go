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

// Message is the entity of email message received
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

// LogProperties is the map to logger attibutes
func (m *Message) LogProperties() map[string]string {
	return map[string]string{
		"message_id":      m.ID,
		"account_id":      m.AccountID,
		"external_id":     m.ExternalID,
		"notification_id": m.NotificationID,
		"extract_id":      m.ExtractID,
		"from":            m.From,
		"to":              m.To,
		"status":          string(m.Status),
		"failure_reason":  m.FailureReason,
		"date":            m.Date.String(),
		"created_at":      m.CreatedAt.String(),
		"updated_at":      m.UpdatedAt.String(),
	}
}

// NewMessage creates a new Message with sensible defaults
func NewMessage(accountID string, from string, to string, externalID string, notificationID string, extractID string, date time.Time) *Message {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

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
