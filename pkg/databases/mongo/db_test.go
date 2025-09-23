package mongo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient_DefaultURI(t *testing.T) {
	c := require.New(t)

	os.Unsetenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := NewClient(ctx)

	c.NoError(err)
	c.NotNil(client)
	c.NotNil(client.client)
}

func TestNewClient_WithEnvURI(t *testing.T) {
	c := require.New(t)

	os.Setenv("MONGO_URI", "mongodb://localhost:27018")
	defer os.Unsetenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	client, err := NewClient(ctx)

	c.NoError(err)
	c.NotNil(client)
}

func TestCollection_ClientNil(t *testing.T) {
	c := require.New(t)

	mc := &MongoClient{client: nil}

	col, err := mc.Collection(TRANSACTIONS, MOVEMENTS)

	c.Nil(col)
	c.Error(err)
	c.Equal("client not initialized", err.Error())
}
