package messageextractor

import (
	"transaction-tracker/internal/movements/domain"
)

type MovementExtractor interface {
	Extract() ([]*domain.Movement, error)
	// SetExtract(*schemas.GmailExtract)
}
