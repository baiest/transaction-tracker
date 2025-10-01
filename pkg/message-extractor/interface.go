package messageextractor

import (
	extractDomain "transaction-tracker/internal/extracts/domain"
	"transaction-tracker/internal/movements/domain"
)

type MovementExtractor interface {
	Extract() ([]*domain.Movement, error)
	SetExtract(*extractDomain.Extract)
}
