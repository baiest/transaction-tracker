package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	_extract_prefix = "EXI"
)

type ExtractStatus string

var (
	ExtractStatusPending   ExtractStatus = "pending"
	ExtractStatusProcessed ExtractStatus = "processed"
	ExtractStatusFailed    ExtractStatus = "failed"
)

type Extract struct {
	ID            string        `bson:"_id,omitempty" json:"id"`
	AccountID     string        `bson:"account_id" json:"account_id"`
	MessageID     string        `bson:"message_id" json:"message_id"`
	InstitutionID string        `bson:"institution_id" json:"institution_id"`
	Month         time.Month    `bson:"month" json:"month"`
	Year          int           `bson:"year" json:"year"`
	Path          string        `bson:"path" json:"path"`
	Status        ExtractStatus `bson:"status" json:"status"`
	CreatedAt     time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at" json:"updated_at"`
}

func NewExtract(accountID string, messageID string, institutionID string, path string, month time.Month, year int) *Extract {
	return &Extract{
		ID:            _extract_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		AccountID:     accountID,
		MessageID:     messageID,
		InstitutionID: institutionID,
		Month:         month,
		Year:          year,
		Path:          path,
		Status:        ExtractStatusPending,
	}
}
