package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"transaction-tracker/internal/accounts/domain"
	"transaction-tracker/internal/accounts/repository"
	"transaction-tracker/pkg/google"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secretKey = "secret"
)

type accountsUseCase struct {
	googleClient *google.GoogleClient
	gmailClient  *google.GmailService
	repo         repository.AccountsRepository
}

func NewAccountsUseCase(googleClient *google.GoogleClient, gmailClient *google.GmailService, repo repository.AccountsRepository) AccountsUseCase {
	return &accountsUseCase{googleClient: googleClient, gmailClient: gmailClient, repo: repo}
}

func (a *accountsUseCase) GetAuthURL() string {
	return a.googleClient.GetAuthURL()
}

func (a *accountsUseCase) CreateAccount(ctx context.Context, account *domain.Account) error {
	return a.repo.CreateAccount(ctx, account)
}

func (a *accountsUseCase) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	return a.repo.GetAccount(ctx, accountID)
}

func (a *accountsUseCase) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	return a.repo.GetAccountByEmail(ctx, email)
}

func (a *accountsUseCase) GetOrCreateAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	account, err := a.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, repository.ErrAccountNotfound) {
			return nil, err
		}

		newAccount := &domain.Account{
			ID:    fmt.Sprintf("acc_%d", time.Now().UnixNano()),
			Email: email,
		}

		if err := a.repo.CreateAccount(ctx, newAccount); err != nil {
			return nil, err
		}

		return newAccount, nil
	}

	return account, err
}

func (a *accountsUseCase) SaveGoogleAccount(ctx context.Context, accountID string, code string) error {
	token, err := a.googleClient.Config.Exchange(ctx, code)
	if err != nil {
		return err
	}

	googleAccount := &google.GoogleAccount{
		Token: token,
	}

	return a.repo.SaveGoogleAccount(ctx, accountID, googleAccount)
}

func (a *accountsUseCase) GenerateTokens(ctx context.Context, account *domain.Account) (string, string, string, error) {
	accessClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    account.ID,
		"email": account.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		"iat":   jwt.NewNumericDate(time.Now()),
	})

	accessToken, err := accessClaims.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", "", err
	}

	refreshClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    account.ID,
		"email": account.Email,
		"exp":   jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iat":   jwt.NewNumericDate(time.Now()),
	})

	refreshToken, err := refreshClaims.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", "", err
	}

	account.RefreshToken = refreshToken

	return accessToken, refreshToken, os.Getenv("REDIRECT_URL"), nil

}

func (a *accountsUseCase) VerifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func (a *accountsUseCase) CreateWatcher(ctx context.Context, account *domain.Account) error {
	projectID := "transaction-tracker-2473"
	topicName := fmt.Sprintf("projects/%s/topics/gmail-notifications", projectID)

	_, _, err := a.gmailClient.CreateWatch(ctx, topicName)
	if err != nil {
		return err
	}

	account.GoogleAccount.IsWatchingGmail = true

	return a.repo.SaveGoogleAccount(ctx, account.ID, account.GoogleAccount)
}

func (a *accountsUseCase) DeleteWatcher(ctx context.Context, account *domain.Account) error {
	err := a.gmailClient.DeleteWatch()
	if err != nil {
		return err
	}

	account.GoogleAccount.IsWatchingGmail = false

	return a.repo.SaveGoogleAccount(ctx, account.ID, account.GoogleAccount)
}
