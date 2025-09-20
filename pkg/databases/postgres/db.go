package postgres

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Client is an interface for our database client
// that the rest of the application will use.
type Client interface {
	GetPool() *pgxpool.Pool
	Close()
}

// client is the concrete implementation of our PostgreSQL client.
type client struct {
	pool *pgxpool.Pool
}

// NewClient creates a new instance of the database client
// with a connection pool.
func NewClient(ctx context.Context) (Client, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Printf("failed to ping database: %v", err)
		pool.Close()
		return nil, fmt.Errorf("database not reachable: %w", err)
	}

	log.Println("Database connection pool established successfully.")
	return &client{pool: pool}, nil
}

// GetPool returns the connection pool.
func (c *client) GetPool() *pgxpool.Pool {
	return c.pool
}

// Close closes the connection pool.
func (c *client) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}
