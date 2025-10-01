package google

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestNewGoogleClient(t *testing.T) {
	c := require.New(t)
	t.Run("Success", func(t *testing.T) {
		clientID = "test_id"
		clientSecret = "test_secret"

		client, err := NewGoogleClient(context.Background())
		c.NoError(err)
		c.NotNil(client)
	})

	t.Run("MissingEnvVars", func(t *testing.T) {
		os.Unsetenv("GOOGLE_CLIENT_ID")
		os.Unsetenv("GOOGLE_CLIENT_SECRET")

		_, err := NewGoogleClient(context.Background())
		c.NoError(err)
	})
}

func TestGetAuthURL(t *testing.T) {
	config := &oauth2.Config{
		ClientID:     "test_id",
		ClientSecret: "test_secret",
		RedirectURL:  "http://localhost",
		Scopes:       []string{"test_scope"},
		Endpoint: oauth2.Endpoint{
			AuthURL: "https://accounts.google.com/o/oauth2/auth",
		},
	}

	client := &GoogleClient{Config: config}

	url := client.GetAuthURL()

	if !strings.Contains(url, "https://accounts.google.com/o/oauth2/auth") {
		t.Errorf("expected URL to contain auth endpoint, but it didn't. Got: %s", url)
	}

	if !strings.Contains(url, "client_id=test_id") {
		t.Errorf("expected URL to contain client ID, but it didn't. Got: %s", url)
	}
}

func TestSetToken(t *testing.T) {
	client := &GoogleClient{}
	token := &oauth2.Token{AccessToken: "test_token"}

	client.SetToken(token)

	if client.token.AccessToken != "test_token" {
		t.Errorf("expected token to be set, but it wasn't")
	}
}

// MockRoundTripper is a mock for http.RoundTripper
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestRefreshToken(t *testing.T) {
	c := require.New(t)

	t.Run("TokenNotExpired", func(t *testing.T) {
		config := &oauth2.Config{}
		client := &GoogleClient{Config: config}
		googleAccount := &GoogleAccount{
			Token: &oauth2.Token{
				Expiry: time.Now().Add(5 * time.Minute),
			},
		}

		_, err := client.RefreshToken(context.Background(), googleAccount)
		c.ErrorContains(err, "failed to refresh token")
	})
}
