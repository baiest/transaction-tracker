package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MovementType defines the type of a movement
type MovementType string

// Source defines the origin of the movement's data based on the processing type.
type Source string

// MovementCategory define category of movements
type MovementCategory string

const (
	_movement_prefix = "MID"

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

	// Salary represents regular job income.
	Salary MovementCategory = "salary"
	// Freelance represents income from freelance or side jobs.
	Freelance MovementCategory = "freelance"
	// Investment represents income from dividends, interests, or capital gains.
	Investment MovementCategory = "investment"

	// Housing represents rent, mortgage, utilities, or related expenses.
	Housing MovementCategory = "housing"
	// Transport represents fuel, public transport, insurance, or maintenance.
	Transport MovementCategory = "transport"
	// Food represents groceries, restaurants, or delivery services.
	Food MovementCategory = "food"
	// Entertainment represents leisure activities, streaming, concerts, or hobbies.
	Entertainment MovementCategory = "entertainment"
	// Shopping represents clothing, electronics, or personal purchases.
	Shopping MovementCategory = "shopping"
	// Health represents medical services, insurance, or medicines.
	Health MovementCategory = "health"
	// Education represents courses, tuition, or books.
	Education MovementCategory = "education"
	// Travel represents vacations or travel-related expenses.
	Travel MovementCategory = "travel"
	// Savings represents money transferred to savings accounts or deposits.
	Savings MovementCategory = "savings"
	// Debt represents loan, credit card, or other debt payments.
	Debt MovementCategory = "debt"

	// Unknown represents an uncategorized movement.
	Unknown MovementCategory = "unknown"
)

// Movement represents a single financial transaction. It's the central business entity.
type Movement struct {
	ID             string           `json:"id" bson:"_id,omitempty"`
	AccountID      string           `json:"account_id" bson:"account_id"`
	InstitutionID  string           `json:"institution_id" bson:"institution_id"`
	MessageID      string           `json:"message_id" bson:"message_id"`
	NotificationID string           `json:"notification_id" bson:"notification_id"`
	Description    string           `json:"description" bson:"description"`
	Amount         float64          `json:"amount" bson:"amount"`
	Type           MovementType     `json:"type" bson:"type"`
	Date           time.Time        `json:"date" bson:"date"`
	Source         Source           `json:"source" bson:"source"`
	Category       MovementCategory `json:"category" bson:"category"`
	CreatedAt      time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" bson:"updated_at"`
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
func NewMovement(accountID string, institutionID string, messageID string, notificationID string, description string, amount float64, category MovementCategory, movementType MovementType, date time.Time, source Source) *Movement {
	return &Movement{
		ID:             _movement_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		AccountID:      accountID,
		InstitutionID:  institutionID,
		MessageID:      messageID,
		NotificationID: notificationID,
		Description:    description,
		Amount:         amount,
		Category:       category,
		Type:           movementType,
		Date:           date,
		Source:         source,
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
