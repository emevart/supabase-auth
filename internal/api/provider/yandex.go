package provider

import (
	"context"
	"strings"

	"github.com/supabase/auth/internal/conf"
	"golang.org/x/oauth2"
)

const (
	defaultYandexAuthBase = "oauth.yandex.ru"
	defaultYandexAPIBase  = "login.yandex.ru"
)

type yandexProvider struct {
	*oauth2.Config
	APIPath string
}

type yandexUser struct {
	ID              string   `json:"id"`
	Login           string   `json:"login"`
	ClientID        string   `json:"client_id"`
	PSUID           string   `json:"psuid"`
	DefaultEmail    string   `json:"default_email"`
	Emails          []string `json:"emails"`
	DefaultAvatarID string   `json:"default_avatar_id"`
	IsAvatarEmpty   bool     `json:"is_avatar_empty"`
	FirstName       string   `json:"first_name"`
	LastName        string   `json:"last_name"`
	DisplayName     string   `json:"display_name"`
	RealName        string   `json:"real_name"`
	Sex             string   `json:"sex"`
	Birthday        string   `json:"birthday"`
}

// NewYandexProvider creates a Yandex ID OAuth2 identity provider.
//
// Yandex ID is a plain OAuth 2.0 provider (not OIDC) used primarily by
// Russian-speaking audiences. Docs: https://yandex.com/dev/id/doc/en/
func NewYandexProvider(ext conf.OAuthProviderConfiguration, scopes string) (OAuthProvider, error) {
	if err := ext.ValidateOAuth(); err != nil {
		return nil, err
	}

	authHost := chooseHost(ext.URL, defaultYandexAuthBase)
	apiHost := chooseHost("", defaultYandexAPIBase)

	// Default Yandex scopes for basic profile. The Yandex authorize endpoint
	// accepts space-separated scopes, which golang.org/x/oauth2 produces by
	// default when joining the Scopes slice.
	oauthScopes := []string{
		"login:info",
		"login:email",
		"login:avatar",
	}

	if scopes != "" {
		oauthScopes = append(oauthScopes, strings.Split(scopes, ",")...)
	}

	return &yandexProvider{
		Config: &oauth2.Config{
			ClientID:     ext.ClientID[0],
			ClientSecret: ext.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:  authHost + "/authorize",
				TokenURL: authHost + "/token",
			},
			Scopes:      oauthScopes,
			RedirectURL: ext.RedirectURI,
		},
		APIPath: apiHost,
	}, nil
}

func (y yandexProvider) GetOAuthToken(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return y.Exchange(ctx, code, opts...)
}

func (y yandexProvider) RequiresPKCE() bool {
	return false
}

func (y yandexProvider) GetUserData(ctx context.Context, tok *oauth2.Token) (*UserProvidedData, error) {
	var u yandexUser
	if err := makeRequest(ctx, tok, y.Config, y.APIPath+"/info?format=json", &u); err != nil {
		return nil, err
	}

	data := &UserProvidedData{}

	// Yandex only exposes verified email addresses: the user cannot attach an
	// unverified address to their Yandex ID account, so we mark everything
	// returned by the /info endpoint as verified.
	if u.DefaultEmail != "" {
		data.Emails = append(data.Emails, Email{
			Email:    u.DefaultEmail,
			Verified: true,
			Primary:  true,
		})
	}

	for _, email := range u.Emails {
		if email == "" || email == u.DefaultEmail {
			continue
		}
		data.Emails = append(data.Emails, Email{
			Email:    email,
			Verified: true,
			Primary:  false,
		})
	}

	var avatarURL string
	if u.DefaultAvatarID != "" && !u.IsAvatarEmpty {
		avatarURL = "https://avatars.yandex.net/get-yapic/" + u.DefaultAvatarID + "/islands-200"
	}

	fullName := u.DisplayName
	if fullName == "" {
		fullName = u.RealName
	}
	if fullName == "" {
		fullName = strings.TrimSpace(u.FirstName + " " + u.LastName)
	}
	if fullName == "" {
		fullName = u.Login
	}

	data.Metadata = &Claims{
		Issuer:            y.APIPath,
		Subject:           u.ID,
		Name:              fullName,
		GivenName:         u.FirstName,
		FamilyName:        u.LastName,
		PreferredUsername: u.Login,
		Picture:           avatarURL,
		Email:             u.DefaultEmail,
		EmailVerified:     u.DefaultEmail != "",

		// To be deprecated, kept for backward compatibility with other providers
		FullName:    fullName,
		AvatarURL:   avatarURL,
		ProviderId:  u.ID,
		UserNameKey: u.Login,
	}

	return data, nil
}
