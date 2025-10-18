package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewAccount(t *testing.T) {
	t.Run("returns an error if email is empty", func(t *testing.T) {
		c := require.New(t)

		account, err := NewAccount("")

		c.Error(err)
		c.Nil(account)
		c.Equal(ErrMissingEmail, err)
	})

	t.Run("returns a new account", func(t *testing.T) {
		c := require.New(t)

		email := "test@example.com"
		account, err := NewAccount(email)

		c.NoError(err)
		c.NotNil(account)
		c.Equal(email, account.Email)
		c.Contains(account.ID, _account_prefix)
	})
}

func TestAccount_LogProperties(t *testing.T) {
	c := require.New(t)

	now := time.Now()
	account := &Account{
		ID:        "test-id",
		Email:     "test@example.com",
		CreatedAt: now,
		UpdatedAt: now,
	}

	props := account.LogProperties()

	c.Equal(account.ID, props["id"])
	c.Equal(account.Email, props["email"])
	c.Equal(account.CreatedAt.String(), props["created_at"])
	c.Equal(account.UpdatedAt.String(), props["updated_at"])
}
