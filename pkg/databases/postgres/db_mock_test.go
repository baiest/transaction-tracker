package postgres

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestNewTestDB(t *testing.T) {
	c := require.New(t)

	db, mock, err := NewTestDB()
	c.NoError(err)
	c.NotNil(db)
	c.Implements((*sqlmock.Sqlmock)(nil), mock)

	sql, err := db.DB()
	c.NoError(err)
	c.NotNil(sql)

	err = sql.Ping()
	c.NoError(err)
}
