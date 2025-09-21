package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"transaction-tracker/pkg/databases"

	"github.com/jackc/pgx/v5/pgxpool"
)

// client is the concrete implementation of our PostgreSQL client.
type client struct {
	pool *pgxpool.Pool
}

var (
	errEnvNotSet = errors.New("DATABASE_URL environment variable is not set")
)

// NewClient creates a new instance of the database client
// with a connection pool.
func NewClient(ctx context.Context) (databases.Client, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errEnvNotSet
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	return &client{pool: pool}, nil
}

// GetPool returns the connection pool.
// nolint:cover
func (c *client) GetPool() *pgxpool.Pool {
	return c.pool
}

// Close closes the connection pool.
// nolint:cover
func (c *client) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}
