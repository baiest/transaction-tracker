package google

import (
	"context"
	"fmt"
	"os"
	"time"

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

// GoogleClient manages the OAuth2 configuration and token for accessing Google APIs.
type GoogleClient struct {
	token  *oauth2.Token
	Config *oauth2.Config
}

// NewGoogleClient creates a new GoogleClient with the necessary OAuth2 configuration.
// It requires GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables to be set.
func NewGoogleClient(ctx context.Context) (*GoogleClient, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GOOGLE_CLIENT_ID y GOOGLE_CLIENT_SECRET must be configurated")
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

	return &GoogleClient{Config: config}, nil
}

// SaveTokenAndInitServices exchanges an authorization code for an OAuth2 token and initializes a GoogleAccount.
func (g *GoogleClient) SaveTokenAndInitServices(ctx context.Context, code string) (*GoogleAccount, error) {
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("error exchange code to token: %w", err)
	}

	g.token = token

	googleAccount := &GoogleAccount{
		Token: token,
	}

	return googleAccount, nil
}

// GetAuthURL returns the URL for the Google authentication page.
func (g *GoogleClient) GetAuthURL() string {
	return g.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

// SetToken sets the OAuth2 token for the GoogleClient.
func (g *GoogleClient) SetToken(token *oauth2.Token) {
	g.token = token
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

// GmailService creates a new GmailService with the provided GoogleAccount credentials.
func (g *GoogleClient) GmailService(ctx context.Context, googleAccount *GoogleAccount) (*GmailService, error) {
	client := g.Config.Client(ctx, googleAccount.Token)

	service, err := NewGmailClient(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("error creating gmail service: %w", err)
	}

	return service, nil
}

// RefreshToken refreshes the OAuth2 token if it has expired or is close to expiring.
func (g *GoogleClient) RefreshToken(ctx context.Context, googleAccount *GoogleAccount) (*oauth2.Token, error) {
	if googleAccount == nil {
		return nil, fmt.Errorf("token was not found")
	}

	expireDate := googleAccount.Token.Expiry

	now := time.Now()
	sub := now.Sub(expireDate)
	if sub > (-3 * time.Minute) {
		return googleAccount.Token, nil
	}

	ts := g.Config.TokenSource(ctx, googleAccount.Token)
	newToken, err := ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	if newToken.AccessToken != googleAccount.Token.AccessToken || newToken.RefreshToken != "" {
		googleAccount.Token = newToken
	}

	g.token = newToken

	return newToken, nil
}
