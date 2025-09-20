package domain

import (
	"errors"
	"time"
)

// MovementType defines the type of a movement
type MovementType string

// Source defines the origin of the movement's data based on the processing type.
type Source string

const (
	// Income represents a money deposit.
	Income MovementType = "income"
	// Expense represents a money withdrawal.
	Expense MovementType = "expense"

	// ManualSource indicates a movement created manually by a user.
	ManualSource Source = "manual"
	// ExtractSource indicates a movement from a processed bank statement.
	ExtractSource Source = "extract"
	// EmailSource indicates a movement from a single email alert.
	EmailSource Source = "email"
)

// Movement represents a single financial transaction. It's the central business entity.
type Movement struct {
	ID            string       `json:"id" bson:"_id,omitempty"`
	AccountID     string       `json:"account_id" bson:"account_id"`
	InstitutionID string       `json:"institution_id" bson:"institution_id"`
	Description   string       `json:"description" bson:"description"`
	Amount        float64      `json:"amount" bson:"amount"`
	Type          MovementType `json:"type" bson:"type"`
	Date          time.Time    `json:"date" bson:"date"`
	Source        Source       `json:"source" bson:"source"`
	Category      string       `json:"category" bson:"category"`
	CreatedAt     time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" bson:"updated_at"`
}

// PaginatedMovements is the structure that encapsulates paginated movements and pagination information.
type PaginatedMovements struct {
	Movements    []*Movement
	TotalRecords int
	TotalPages   int
	CurrentPage  int
	Limit        int
	Offset       int
}

// NewMovement creates an instece of Movement.
func NewMovement(accountID string, institutionID string, description string, amount float64, movementType MovementType, date time.Time, source Source) *Movement {
	return &Movement{
		AccountID:     accountID,
		InstitutionID: institutionID,
		Description:   description,
		Amount:        amount,
		Type:          movementType,
		Date:          date,
		Source:        source,
	}
}

// ParseMovementType validates a string and returns the corresponding MovementType.
// It returns an error if the string does not match a known movement type (Income or Expense).
func ParseMovementType(t string) (MovementType, error) {
	switch MovementType(t) {
	case Income, Expense:
		return MovementType(t), nil
	default:
		return "", errors.New("invalid movement type")
	}
}
