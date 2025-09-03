package schemas

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Movement struct {
	ID         string    `bson:"_id"`
	Email      string    `bson:"email"`
	ExtractID  string    `bson:"extract_id"`
	Date       time.Time `bson:"date"`
	Value      float64   `bson:"value"`
	IsNegative bool      `bson:"is_negative"`
	Type       string    `bson:"type"`
	Detail     string    `bson:"detail"`
}

const _movement_prefix = "MOV"

func NewMovement(email string, extractID string, date time.Time, value float64, isNegative bool, movementType string, detail string) *Movement {
	return &Movement{
		ID:         _movement_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email:      email,
		ExtractID:  extractID,
		Date:       date,
		Value:      value,
		IsNegative: isNegative,
		Type:       movementType,
		Detail:     detail,
	}
}
