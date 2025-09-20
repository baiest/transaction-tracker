package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewClient_Failed(t *testing.T) {
	c := require.New(t)
	os.Unsetenv("DATABASE_URL")

	_, err := NewClient(context.Background())
	c.ErrorIs(err, errEnvNotSet)
}

func TestNewClient_Success(t *testing.T) {
	c := require.New(t)
	os.Setenv("DATABASE_URL", "dummy")

	_, err := NewClient(context.Background())
	c.Error(err)

	os.Setenv("DATABASE_URL", "postgres://dummy")

	_, err = NewClient(context.Background())
	c.Error(err)
}

func TestClient_GetPool_And_Close(t *testing.T) {
	c := require.New(t)

	cl := &client{pool: nil}

	c.Nil(cl.GetPool())

	cl.Close()
}
