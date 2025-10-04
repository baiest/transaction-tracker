package usecase

import (
	"context"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/extracts/domain"
)

type ExtractsUsecase interface {
	GetByMessageID(ctx context.Context, messageID string) (*domain.Extract, error)
	GetExtractMessages(ctx context.Context, bankName string, account *accountsDomain.Account) ([]string, error)
	Save(ctx context.Context, extract *domain.Extract) error
	Update(ctx context.Context, extract *domain.Extract) error
}
