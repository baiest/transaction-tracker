package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"

	"transaction-tracker/internal/movements/domain"
)

var (
	fixedTime = time.Date(2025, 9, 20, 12, 0, 0, 0, time.UTC)
)

// NewTestPostgresRepository is a constructor exclusively for tests.
// It accepts the mock interface and returns the repository implementation.
func NewTestPostgresRepository(db pgxmock.PgxPoolIface) MovementRepository {
	return &postgresRepository{db: db, nowFunc: func() time.Time { return fixedTime }}
}

func setupMockDB(t *testing.T) (MovementRepository, pgxmock.PgxPoolIface, func()) {
	mockPool, err := pgxmock.NewPool()
	require.NoError(t, err)

	repo := NewTestPostgresRepository(mockPool)

	cleanup := func() { mockPool.Close() }

	return repo, mockPool, cleanup
}

func TestCreateMovement(t *testing.T) {
	c := require.New(t)

	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	movement := &domain.Movement{
		ID:            uuid.New().String(),
		AccountID:     "acc1",
		InstitutionID: "inst1",
		Description:   "Test Description",
		Amount:        1000.0,
		Type:          "expense",
		Date:          now,
		Source:        "card",
		Category:      "groceries",
		CreatedAt:     fixedTime,
		UpdatedAt:     fixedTime,
	}

	mock.ExpectExec(`INSERT INTO movements`).
		WithArgs(
			movement.ID,
			movement.AccountID,
			movement.InstitutionID,
			movement.Description,
			movement.Amount,
			movement.Type,
			movement.Date,
			movement.Source,
			movement.Category,
			movement.CreatedAt,
			movement.UpdatedAt,
		).WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err := repo.CreateMovement(context.Background(), movement)
	c.NoError(err)

	c.NoError(mock.ExpectationsWereMet())
}

func TestGetMovementByID(t *testing.T) {
	c := require.New(t)

	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()
	instID := "inst1"
	desc := "Test Desc"
	amount := float64(1000.0)
	date := now
	source := "card"
	cat := "groceries"

	columns := []string{"id", "account_id", "institution_id", "description", "amount", "type", "date", "source", "category", "created_at", "updated_at"}
	rows := pgxmock.NewRows(columns).
		AddRow("mov1", "acc1", &instID, &desc, amount, "expense", &date, &source, &cat, &now, &now)

	mock.ExpectQuery(`SELECT (.+) FROM movements WHERE id = \$1`).
		WithArgs("mov1").
		WillReturnRows(rows)

	m, err := repo.GetMovementByID(context.Background(), "mov1")
	c.NoError(err)
	c.Equal("mov1", m.ID)
	c.Equal("acc1", m.AccountID)
	c.Equal(1000.0, m.Amount)
	c.NoError(mock.ExpectationsWereMet())
}

func TestGetMovementsByAccountID(t *testing.T) {
	c := require.New(t)

	repo, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	instID1 := "inst1"
	desc1 := "Desc 1"
	amount1 := 1000.0
	date1 := now
	source1 := "card"
	cat1 := "groceries"

	instID2 := "inst1"
	desc2 := "Desc 2"
	amount2 := 2000.0
	date2 := now.Add(time.Hour)
	source2 := "transfer"
	cat2 := "salary"

	columns := []string{"id", "account_id", "institution_id", "description", "amount", "type", "date", "source", "category", "created_at", "updated_at"}
	rows := pgxmock.NewRows(columns).
		AddRow("mov1", "acc1", &instID1, &desc1, amount1, "expense", &date1, &source1, &cat1, &now, &now).
		AddRow("mov2", "acc1", &instID2, &desc2, amount2, "income", &date2, &source2, &cat2, &now, &now)

	mock.ExpectQuery(`SELECT (.+) FROM movements WHERE account_id = \$1.*LIMIT \$2 OFFSET \$3`).
		WithArgs("acc1", 10, 1).
		WillReturnRows(rows)

	movements, err := repo.GetMovementsByAccountID(context.Background(), "acc1", 10, 1)
	c.NoError(err)

	c.Len(movements, 2)
	c.Equal("mov1", movements[0].ID)
	c.Equal(2000.0, movements[1].Amount)
	c.NoError(mock.ExpectationsWereMet())
}
