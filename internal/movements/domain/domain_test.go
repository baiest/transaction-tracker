package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewMovement(t *testing.T) {
	c := require.New(t)

	testAccountID := "acc-123"
	testInstitutionID := "inst-456"
	testMessageID := "mid123"
	testNotificationID := "nid1"
	testDescription := "Test transaction"
	testAmount := 125.75
	testMovementType := Income
	testDate := time.Now().UTC()
	testSource := ManualSource
	testCategory := Unknown

	movement := NewMovement(
		testAccountID,
		testInstitutionID,
		testMessageID,
		testNotificationID,
		testDescription,
		testAmount,
		testCategory,
		testMovementType,
		testDate,
		testSource,
	)

	c.Equal(testAccountID, movement.AccountID)
	c.Equal(testInstitutionID, movement.InstitutionID)
	c.Equal(testDescription, movement.Description)
	c.Equal(testAmount, movement.Amount)
	c.Equal(testCategory, movement.Category)
	c.Equal(testMovementType, movement.Type)
	c.Equal(testSource, movement.Source)

	c.NotEmpty(movement.ID)
	c.True(movement.CreatedAt.IsZero())
	c.True(movement.UpdatedAt.IsZero())
}

func TestParseMovementType(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      MovementType
		wantError string
	}{
		{
			name:      "valid Income",
			input:     "income",
			want:      Income,
			wantError: "",
		},
		{
			name:      "valid Expense",
			input:     "expense",
			want:      Expense,
			wantError: "",
		},
		{
			name:      "invalid type",
			input:     "Saving",
			want:      "",
			wantError: "invalid movement type",
		},
		{
			name:      "empty string",
			input:     "",
			want:      "",
			wantError: "invalid movement type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := require.New(t)

			got, err := ParseMovementType(tt.input)

			if tt.wantError != "" {
				c.Error(err)
				c.EqualError(err, tt.wantError)
				c.Empty(got)
			} else {
				c.NoError(err)
				c.Equal(tt.want, got)
			}
		})
	}
}
