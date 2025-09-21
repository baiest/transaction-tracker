package databases

import "github.com/jackc/pgx/v5/pgxpool"

// Client is an interface for our database client
// that the rest of the application will use.
type Client interface {
	GetPool() *pgxpool.Pool
	Close()
}
