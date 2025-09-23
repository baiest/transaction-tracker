package repository

import (
	"context"
	"time"
	"transaction-tracker/internal/movements/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBQuerier is the interface that abstracts the database methods we need.
type DBQuerier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type postgresRepository struct {
	db      DBQuerier
	nowFunc func() time.Time
}

func NewPostgresRepository(db *pgxpool.Pool) MovementRepository {
	return &postgresRepository{db: db, nowFunc: time.Now}
}

// CreateMovement saves a movement using database/sql.
func (r *postgresRepository) CreateMovement(ctx context.Context, movement *domain.Movement) error {
	now := r.nowFunc()

	movement.CreatedAt = now
	movement.UpdatedAt = now

	query := `INSERT INTO movements (
	id,
	account_id,
	institution_id,
	message_id,
	notification_id,
	description,
	amount,
	type,
	date,
	source,
	category,
	created_at,
	updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := r.db.Exec(ctx, query,
		movement.ID,
		movement.AccountID,
		movement.InstitutionID,
		movement.MessageID,
		movement.NotificationID,
		movement.Description,
		movement.Amount,
		movement.Type,
		movement.Date,
		movement.Source,
		movement.Category,
		movement.CreatedAt,
		movement.UpdatedAt)

	return err
}

// GetMovementByID gets a movement by ID.
func (r *postgresRepository) GetMovementByID(ctx context.Context, id string, accountID string) (*domain.Movement, error) {
	query := `SELECT
	id, account_id, institution_id, message_id, notification_id, description, amount, type, date, source, category, created_at, updated_at
	FROM movements
	WHERE id = $1 AND account_id = $2`

	row := r.db.QueryRow(ctx, query, id, accountID)

	return scanToMovement(row.Scan)
}

// GetTotalMovementsByAccountID retrieves the total count of movements for a given account.
func (r *postgresRepository) GetTotalMovementsByAccountID(ctx context.Context, accountID string) (int, error) {
	countQuery := `SELECT COUNT(*) FROM movements WHERE account_id = $1`
	var totalRecords int
	err := r.db.QueryRow(ctx, countQuery, accountID).Scan(&totalRecords)
	if err != nil {
		return 0, err
	}

	return totalRecords, nil
}

// GetMovementsByAccountID gets a user's movements with pagination.
func (r *postgresRepository) GetMovementsByAccountID(ctx context.Context, accountID string, limit int, offset int) ([]*domain.Movement, error) {
	query := `SELECT
		id, account_id, institution_id, message_id, notification_id, description, amount, type, date, source, category, created_at, updated_at
	FROM movements
	WHERE account_id = $1
	ORDER BY date DESC
	LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, accountID, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movements []*domain.Movement
	for rows.Next() {
		m, err := scanToMovement(rows.Scan)
		if err != nil {
			return nil, err
		}

		movements = append(movements, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movements, nil
}

func scanToMovement(scanFn func(...any) error) (*domain.Movement, error) {
	m := &domain.Movement{}

	var institutionID, messageID, notificationID, description, source *string
	var date, createdAt, updatedAt *time.Time
	var movementType, category string

	err := scanFn(
		&m.ID, &m.AccountID, &institutionID, &messageID, &notificationID, &description, &m.Amount,
		&movementType, &date, &source, &category, &createdAt, &updatedAt,
	)

	if err != nil {
		return nil, err
	}

	if institutionID != nil {
		m.InstitutionID = *institutionID
	}

	if messageID != nil {
		m.MessageID = *messageID
	}

	if notificationID != nil {
		m.NotificationID = *notificationID
	}

	if description != nil {
		m.Description = *description
	}

	if date != nil {
		m.Date = *date
	}

	if source != nil {
		m.Source = domain.Source(*source)
	}

	m.Type = domain.MovementType(movementType)
	m.Category = domain.MovementCategory(category)

	if createdAt != nil {
		m.CreatedAt = *createdAt
	}

	if updatedAt != nil {
		m.UpdatedAt = *updatedAt
	}

	return m, nil
}
