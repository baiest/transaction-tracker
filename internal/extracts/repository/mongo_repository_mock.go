package repository

import (
	"context"
	"transaction-tracker/internal/extracts/domain"

	"github.com/stretchr/testify/mock"
)

// MockExtractssRepository is a mock of ExtractsRepository
type MockExtractsRepository struct {
	mock.Mock
}

func (r *MockExtractsRepository) GetByMessageID(ctx context.Context, messageID string) (*domain.Extract, error) {
	args := r.Called(ctx, messageID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Extract), args.Error(1)
}

func (r *MockExtractsRepository) Save(ctx context.Context, extract *domain.Extract) error {
	args := r.Called(ctx, extract)
	return args.Error(0)
}

func (r *MockExtractsRepository) Update(ctx context.Context, extract *domain.Extract) error {
	args := r.Called(ctx, extract)
	return args.Error(0)
}
