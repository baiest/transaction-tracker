package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"transaction-tracker/pkg/databases"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// client is the concrete implementation of our PostgreSQL client.
type client struct {
	pool *pgxpool.Pool
}

type queryTracer struct{}

func (qt *queryTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	fmt.Printf("[QUERY START] %s | args=%v\n", data.SQL, data.Args)
	return ctx
}

func (qt *queryTracer) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	if data.Err != nil {
		fmt.Printf("[QUERY END] ERROR: %v\n", data.Err)
	} else {
		fmt.Printf("[QUERY END] commandTag=%v\n", data.CommandTag.String())
	}
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

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Activar logs
	config.ConnConfig.Tracer = &queryTracer{}

	pool, err := pgxpool.NewWithConfig(ctx, config)
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
