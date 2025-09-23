package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"transaction-tracker/logger"
	"transaction-tracker/pkg/databases"

	loggerModels "transaction-tracker/logger/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// client is the concrete implementation of our PostgreSQL client.
type client struct {
	pool *pgxpool.Pool
}

type queryTracer struct {
	log *loggerModels.Logger
}

func (qt *queryTracer) TraceQueryStart(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryStartData,
) context.Context {
	return ctx
}

func (qt *queryTracer) TraceQueryEnd(
	ctx context.Context,
	conn *pgx.Conn,
	data pgx.TraceQueryEndData,
) {
	if data.Err != nil {
		qt.log.Error(loggerModels.LogProperties{
			Event: "postgres_query_failed",
			Error: data.Err,
			AdditionalParams: []loggerModels.Properties{
				logger.MapToProperties(map[string]string{
					"query": data.CommandTag.String(),
				}),
			},
		})
	}
}

var (
	errEnvNotSet = errors.New("DATABASE_URL environment variable is not set")
)

// NewClient creates a new instance of the database client
// with a connection pool.
func NewClient(ctx context.Context) (databases.Client, error) {
	l, err := logger.GetLogger(ctx, "postgres")
	if err != nil {
		return nil, err
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, errEnvNotSet
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Activar logs
	config.ConnConfig.Tracer = &queryTracer{log: l}

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
