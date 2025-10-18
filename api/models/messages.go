package models

import (
	"time"
	"transaction-tracker/internal/messages/domain"
)

type MessageResponse struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	ExternalID     string    `json:"external_id,omitempty"`
	NotificationID string    `json:"notification_id,omitempty"`
	ExtractID      string    `json:"extract_id,omitempty"`
	From           string    `json:"from,omitempty"`
	To             string    `json:"to,omitempty"`
	Status         string    `json:"status"`
	Date           time.Time `json:"date"`
}

type MessageRequest struct {
	ExternalID string `form:"external_id" binding:"required"`
}

// ToMovementResponse converts a single domain.Movement to an API MovementResponse.
func ToMessageResponse(m *domain.Message) *MessageResponse {
	return &MessageResponse{
		ID:             m.ID,
		AccountID:      m.AccountID,
		ExternalID:     m.ExternalID,
		NotificationID: m.NotificationID,
		ExtractID:      "",
		From:           m.From,
		To:             m.To,
		Status:         string(m.Status),
		Date:           m.Date,
	}
}
