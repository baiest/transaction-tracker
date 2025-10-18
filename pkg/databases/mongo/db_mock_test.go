package mongo

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMockMongoClient(t *testing.T) {
	c := require.New(t)

	mock := NewMockMongoClient()

	col, err := mock.Collection(TRANSACTIONS, MOVEMENTS)
	c.NoError(err)
	c.NotNil(col)

	// Should return the same pointer for the same db/collection
	col2, err := mock.Collection(TRANSACTIONS, MOVEMENTS)
	c.NoError(err)
	c.Equal(col, col2)

	// Force an error
	mock.Fail = true
	col3, err := mock.Collection(TRANSACTIONS, ACCOUNTS)
	c.Error(err)
	c.Nil(col3)
}
