package google

import (
	"context"
	"fmt"
	"os"
	"transaction-tracker/database/mongo/schemas"
	"transaction-tracker/googleapi/repositories"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var (
	clientID     = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL  = os.Getenv("GOOGLE_REDIRECT_URL")
)

type GoogleClient struct {
	token        *oauth2.Token
	email        string
	Config       *oauth2.Config
	gmailService *GmailService
	repository   *repositories.GoogleAccountsRepository
}

func NewGoogleClient(ctx context.Context) (*GoogleClient, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID y GOOGLE_CLIENT_SECRET must be configurated")
	}

	repository, err := repositories.NeGoogleAccountsRepository(ctx)
	if err != nil {
		return nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			gmail.GmailReadonlyScope,
			gmail.GmailModifyScope,
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleClient{Config: config, repository: repository, email: ""}, nil
}

func (g *GoogleClient) SaveTokenAndInitServices(ctx context.Context, code string) error {
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("Error exchange code to toke: %s", err)
	}

	g.token = token

	email, err := g.GetUserEmail(ctx)
	if err != nil {
		return err
	}

	g.email = email

	googleAccount := &schemas.GoogleAccount{
		ID:    email,
		Token: token,
	}

	err = g.repository.SaveToken(ctx, googleAccount)
	if err != nil {
		return err
	}

	service, err := NewGmailClient(ctx, g)
	if err != nil {
		return fmt.Errorf("Error creating gmail service: %v", err)
	}

	g.gmailService = service

	return nil
}

func (g *GoogleClient) GetAuthURL() string {
	return g.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

// GetUserEmail fetches the authenticated user's email address using the Gmail API profile endpoint.
func (g *GoogleClient) GetUserEmail(ctx context.Context) (string, error) {
	if g.token == nil {
		return "", fmt.Errorf("missing oauth token")
	}

	httpClient := g.Config.Client(ctx, g.token)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", fmt.Errorf("error creating gmail service: %w", err)
	}

	profile, err := service.Users.GetProfile("me").Do()
	if err != nil {
		return "", fmt.Errorf("error getting gmail profile: %w", err)
	}

	return profile.EmailAddress, nil
}

func (g *GoogleClient) GmailService(ctx context.Context) (*GmailService, error) {
	if g.email != "" {
		account, err := g.repository.GetTokenByEmail(ctx, g.email)
		if err != nil {
			return nil, err
		}

		g.token = account.Token

		service, err := NewGmailClient(ctx, g)
		if err != nil {
			return nil, fmt.Errorf("Error creating gmail service: %v", err)
		}

		g.gmailService = service
	}

	if g.gmailService == nil {
		return nil, fmt.Errorf("gmail service not inicialized")
	}

	return g.gmailService, nil
}

func (g *GoogleClient) SetEmail(email string) {
	g.email = email
}

func (g *GoogleClient) Email() string {
	return g.email
}

func (g *GoogleClient) RefreshToken(ctx context.Context) (*oauth2.Token, error) {
	if g.email == "" {
		return nil, fmt.Errorf("email not set in GoogleClient")
	}

	account, err := g.repository.GetTokenByEmail(ctx, g.email)
	if err != nil {
		return nil, fmt.Errorf("failed to get token from repository: %w", err)
	}

	if account.Token == nil {
		return nil, fmt.Errorf("no token found for email: %s", g.email)
	}

	ts := g.Config.TokenSource(ctx, account.Token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.AccessToken != account.Token.AccessToken || newToken.RefreshToken != "" {
		account.Token = newToken
		if err := g.repository.SaveToken(ctx, account); err != nil {
			return nil, fmt.Errorf("failed to save refreshed token: %w", err)
		}
	}

	g.token = newToken

	return newToken, nil
}
