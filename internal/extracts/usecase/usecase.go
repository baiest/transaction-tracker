package usecase

import (
	"context"
	"errors"

	accountsDomain "transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/extracts/domain"
	"transaction-tracker/internal/extracts/repository"
	"transaction-tracker/pkg/google"
)

var (
	// ErrExtractNotFound is returned when an extract is not found in the repository.
	ErrExtractNotFound = repository.ErrExtractNotFound
)

type extractsUsecase struct {
	repo         repository.ExtractsRepository
	googleClient google.GoogleClientAPI
}

// NewExtractsUsecase creates a new instance of ExtractsUsecase.
func NewExtractsUsecase(googleClient google.GoogleClientAPI, repo repository.ExtractsRepository) ExtractsUsecase {
	return &extractsUsecase{
		repo:         repo,
		googleClient: googleClient,
	}
}

// GetByMessageID retrieves an extract by its message ID.
func (u *extractsUsecase) GetByMessageID(ctx context.Context, messageID string) (*domain.Extract, error) {
	extract, err := u.repo.GetByMessageID(ctx, messageID)
	if err != nil {
		if errors.Is(err, repository.ErrExtractNotFound) {
			return nil, ErrExtractNotFound
		}
	}

	return extract, err
}

func (u *extractsUsecase) GetExtractMessages(ctx context.Context, bankName string, account *accountsDomain.Account) ([]string, error) {
	gmailService, err := u.googleClient.GmailService(ctx, account.GoogleAccount)
	if err != nil {
		return nil, err
	}

	response, err := gmailService.GetExtractMessages(ctx, bankName)
	if err != nil {
		return nil, err
	}

	messageIDs := []string{}

	for _, msg := range response.Messages {
		messageIDs = append(messageIDs, msg.Id)
	}

	return messageIDs, nil
}

// Save saves the given extract using the repository.
func (u *extractsUsecase) Save(ctx context.Context, extract *domain.Extract) error {
	return u.repo.Save(ctx, extract)
}

// Update updates the given extract.
// This is a placeholder implementation and should be replaced with actual logic.
func (u *extractsUsecase) Update(ctx context.Context, extract *domain.Extract) error {
	return u.repo.Update(ctx, extract)
}
