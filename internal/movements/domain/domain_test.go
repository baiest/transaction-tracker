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
	testDescription := "Test transaction"
	testAmount := 125.75
	testMovementType := Income
	testDate := time.Now().UTC()
	testSource := ManualSource

	movement := NewMovement(
		testAccountID,
		testInstitutionID,
		testDescription,
		testAmount,
		testMovementType,
		testDate,
		testSource,
	)

	c.Equal(testAccountID, movement.AccountID)
	c.Equal(testInstitutionID, movement.InstitutionID)
	c.Equal(testDescription, movement.Description)
	c.Equal(testAmount, movement.Amount)
	c.Equal(testMovementType, movement.Type)
	c.Equal(testSource, movement.Source)

	c.Empty(movement.ID)
	c.True(movement.CreatedAt.IsZero())
	c.True(movement.UpdatedAt.IsZero())
}
