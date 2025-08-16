package googleapi

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var (
	clientID     = os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL  = os.Getenv("GOOGLE_REDIRECT_URL")
)

type GoogleClient struct {
	token        *oauth2.Token
	Config       *oauth2.Config
	gmailService *GmailService
}

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

func (g *GoogleClient) SaveTokenAndInitServices(ctx context.Context, code string) error {
	token, err := g.Config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("Error exchange code to toke: %s", err)
	}

	g.token = token

	service, err := NewGmailService(ctx, g)
	if err != nil {
		return fmt.Errorf("Error creating gmail service: %v", err)
	}

	g.gmailService = service

	return nil
}

func (g *GoogleClient) GetAuthURL() string {
	return g.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (g *GoogleClient) GmailService() (*GmailService, error) {
	if g.gmailService == nil {
		return nil, fmt.Errorf("gmail service not inicialized")
	}

	return g.gmailService, nil
}
