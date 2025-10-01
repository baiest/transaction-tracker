package domain

import (
	"errors"
	"strings"
	"time"
	"transaction-tracker/pkg/google"

	"github.com/google/uuid"
)

var (
	_account_prefix = "AID"
	secretKey       = "secret"

	// ErrMissingEmail is returned when an email is not provided to create a new account.
	ErrMissingEmail = errors.New("missing email")
)

// Account represents a user account in the system.
type Account struct {
	ID            string                `bson:"_id"`
	Email         string                `bson:"email"`
	RefreshToken  string                `bson:"refresh_token"`
	GoogleAccount *google.GoogleAccount `bson:"google_account"`
	CreatedAt     time.Time             `bson:"created_at"`
	UpdatedAt     time.Time             `bson:"updated_at"`
}

// LogProperties returns a map of the account's properties for logging.
func (a *Account) LogProperties() map[string]string {
	return map[string]string{
		"id":         a.ID,
		"email":      a.Email,
		"created_at": a.CreatedAt.String(),
		"updated_at": a.UpdatedAt.String(),
	}
}

// NewAccount creates a new Account with the provided email.
func NewAccount(email string) (*Account, error) {
	if email == "" {
		return nil, ErrMissingEmail
	}

	return &Account{
		ID:    _account_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email: email,
	}, nil
}
