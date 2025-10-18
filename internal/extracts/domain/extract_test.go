package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewExtract_SetsFields(t *testing.T) {
	c := require.New(t)

	accountID := "acc-123"
	messageID := "msg-456"
	institutionID := "inst-789"
	path := "/some/path.pdf"
	month := time.January
	year := 2025

	ex := NewExtract(accountID, messageID, institutionID, path, month, year)

	c.NotNil(ex)
	c.True(strings.HasPrefix(ex.ID, _extract_prefix))
	c.NotContains(ex.ID, "-") // uuid dashes removed
	c.Equal(accountID, ex.AccountID)
	c.Equal(messageID, ex.MessageID)
	c.Equal(institutionID, ex.InstitutionID)
	c.Equal(path, ex.Path)
	c.Equal(month, ex.Month)
	c.Equal(year, ex.Year)
	c.Equal(ExtractStatusPending, ex.Status)

	// CreatedAt and UpdatedAt are zero values on construction
	c.True(ex.CreatedAt.IsZero())
	c.True(ex.UpdatedAt.IsZero())
}

func TestNewExtract_UniqueIDs(t *testing.T) {
	c := require.New(t)

	ex1 := NewExtract("a", "m1", "i", "p", time.February, 2024)
	ex2 := NewExtract("a", "m2", "i", "p", time.February, 2024)

	c.NotEqual(ex1.ID, ex2.ID)
}
