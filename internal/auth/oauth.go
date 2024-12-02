package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goBookMarker/internal/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

type OAuthProvider struct {
	config *oauth2.Config
	token  *oauth2.Token
}

type UserInfo struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Provider string `json:"provider"`
}

func NewGoogleAuth() *OAuthProvider {
	return &OAuthProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  "com.gobookmarker:/oauth2callback",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/drive.appdata",
			},
			Endpoint: google.Endpoint,
		},
	}
}

func NewMicrosoftAuth() *OAuthProvider {
	return &OAuthProvider{
		config: &oauth2.Config{
			ClientID:     os.Getenv("MS_CLIENT_ID"),
			ClientSecret: os.Getenv("MS_CLIENT_SECRET"),
			RedirectURL:  "com.gobookmarker:/oauth2callback",
			Scopes: []string{
				"offline_access",
				"User.Read",
				"Files.ReadWrite.AppFolder",
			},
			Endpoint: microsoft.AzureADEndpoint("common"),
		},
	}
}

func (p *OAuthProvider) GetAuthURL() string {
	return p.config.AuthCodeURL("state")
}

func (p *OAuthProvider) HandleCallback(code string) error {
	ctx := context.Background()
	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("token exchange error: %v", err)
	}
	p.token = token
	return nil
}

func (p *OAuthProvider) GetUserInfo() (*models.User, error) {
	if p.token == nil {
		return nil, fmt.Errorf("no token available")
	}

	client := p.config.Client(context.Background(), p.token)

	// Different endpoints for different providers
	var userInfoURL string
	switch p.config.Endpoint {
	case google.Endpoint:
		userInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
	case microsoft.AzureADEndpoint("common"):
		userInfoURL = "https://graph.microsoft.com/v1.0/me"
	default:
		return nil, fmt.Errorf("unsupported OAuth provider")
	}

	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &models.User{
		Email: userInfo.Email,
		Name:  userInfo.Name,
	}, nil
}

func (p *OAuthProvider) GetToken() *oauth2.Token {
	return p.token
}

func (p *OAuthProvider) SetToken(token *oauth2.Token) {
	p.token = token
}
