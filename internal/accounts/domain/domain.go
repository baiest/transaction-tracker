package domain

import (
	"errors"
	"strings"
	"time"
	"transaction-tracker/pkg/google"

	"github.com/google/uuid"
)

type Account struct {
	ID            string                `bson:"_id"`
	Email         string                `bson:"email"`
	RefreshToken  string                `bson:"refresh_token"`
	GoogleAccount *google.GoogleAccount `bson:"google_account"`
	CreatedAt     time.Time             `bson:"created_at"`
	UpdatedAt     time.Time             `bson:"updated_at"`
}

var (
	_account_prefix = "AID"
	secretKey       = "secret"

	ErrMissingEmail = errors.New("missing email")
)

func NewAccount(email string) (*Account, error) {
	if email == "" {
		return nil, ErrMissingEmail
	}

	return &Account{
		ID:    _account_prefix + strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email: email,
	}, nil
}
