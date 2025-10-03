package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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
	googleClient google.GoogleClientAPI
	repo         repository.AccountsRepository
}

// NewAccountsUseCase creates a new AccountsUsecase with the provided Google client and repository.
func NewAccountsUseCase(googleClient google.GoogleClientAPI, repo repository.AccountsRepository) AccountsUsecase {
	return &accountsUseCase{googleClient: googleClient, repo: repo}
}

// GetAuthURL returns the URL for the Google authentication page.
func (a *accountsUseCase) GetAuthURL() string {
	return a.googleClient.GetAuthURL()
}

// CreateAccount creates a new account.
func (a *accountsUseCase) CreateAccount(ctx context.Context, account *domain.Account) error {
	return a.repo.CreateAccount(ctx, account)
}

// GetAccount retrieves an account by its ID.
func (a *accountsUseCase) GetAccount(ctx context.Context, accountID string) (*domain.Account, error) {
	return a.repo.GetAccount(ctx, accountID)
}

// GetAccountByEmail retrieves an account by its email address.
func (a *accountsUseCase) GetAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	email = strings.ToLower(email)

	return a.repo.GetAccountByEmail(ctx, email)
}

// GetOrCreateAccountByEmail retrieves an account by its email address, or creates a new one if it doesn't exist.
func (a *accountsUseCase) GetOrCreateAccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	account, err := a.repo.GetAccountByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, repository.ErrAccountNotfound) {
			return nil, err
		}

		newAccount, err := domain.NewAccount(email)
		if err != nil {
			return nil, err
		}

		if err := a.repo.CreateAccount(ctx, newAccount); err != nil {
			return nil, err
		}

		return newAccount, nil
	}

	return account, err
}

// SaveGoogleAccount exchanges the authorization code for a token, retrieves the user's email, and saves the Google account information.
func (a *accountsUseCase) SaveGoogleAccount(ctx context.Context, code string) (*domain.Account, error) {
	token, err := a.googleClient.Config().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	googleAccount := &google.GoogleAccount{
		Token: token,
	}

	a.googleClient.SetToken(token)

	email, err := a.googleClient.GetUserEmail(ctx)
	if err != nil {
		return nil, err
	}

	account, err := a.GetOrCreateAccountByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	account.GoogleAccount = googleAccount

	return account, a.repo.SaveGoogleAccount(ctx, account.ID, googleAccount)
}

// GenerateTokens generates new access and refresh tokens for the provided account.
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

	err = a.RefreshGoogleToken(ctx, account)
	if err != nil {
		return "", "", "", err
	}

	err = a.repo.UpdateAccount(ctx, account)
	if err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, os.Getenv("REDIRECT_URL"), nil
}

// RefreshGoogleToken refreshes the Google OAuth2 token for the provided account.
func (a *accountsUseCase) RefreshGoogleToken(ctx context.Context, account *domain.Account) error {
	_, err := a.googleClient.RefreshToken(ctx, account.GoogleAccount)
	if err != nil {
		return err
	}

	return a.repo.SaveGoogleAccount(ctx, account.ID, account.GoogleAccount)
}

// VerifyToken verifies the given JWT token string.
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

// CreateWatcher creates a new Gmail watcher for the provided account.
func (a *accountsUseCase) CreateWatcher(ctx context.Context, account *domain.Account) error {
	projectID := "transaction-tracker-2473"
	topicName := fmt.Sprintf("projects/%s/topics/gmail-notifications", projectID)

	gmailClient, err := google.NewGmailClient(ctx, a.googleClient.Config().Client(ctx, account.GoogleAccount.Token))
	if err != nil {
		return err
	}

	_, _, err = gmailClient.CreateWatch(ctx, topicName)
	if err != nil {
		return err
	}

	account.GoogleAccount.IsWatchingGmail = true

	return a.repo.SaveGoogleAccount(ctx, account.ID, account.GoogleAccount)
}

// DeleteWatcher deletes the Gmail watcher for the provided account.
func (a *accountsUseCase) DeleteWatcher(ctx context.Context, account *domain.Account) error {
	gmailClient, err := google.NewGmailClient(ctx, a.googleClient.Config().Client(ctx, account.GoogleAccount.Token))
	if err != nil {
		return err
	}

	err = gmailClient.DeleteWatch()
	if err != nil {
		return err
	}

	account.GoogleAccount.IsWatchingGmail = false

	return a.repo.SaveGoogleAccount(ctx, account.ID, account.GoogleAccount)
}

// UpdateAccount updates the provided account.
func (a *accountsUseCase) UpdateAccount(ctx context.Context, account *domain.Account) error {
	return a.repo.UpdateAccount(ctx, account)
}
