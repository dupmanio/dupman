package credentials

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// ClientCredentialsProvider implements Credentials Provider
// with OAuth client_credentials Grant Type.
type ClientCredentialsProvider struct {
	config clientcredentials.Config
	token  *oauth2.Token
}

// NewClientCredentials creates a new instance of the ClientCredentialsProvider.
//
// Example:
//
//	// Create new instance of Credential Provider, e.g ClientCredentials.
//	cred, err := credentials.NewClientCredentials("worker", "some-secret-key", []string{"system"})
//
//	// Create new session.
//	sess, err := session.New(&dupman.Config{Credentials: cred})
func NewClientCredentials(clientID, clientSecret string, scopes []string) (*ClientCredentialsProvider, error) {
	ctx := context.Background()
	// @todo: update url!
	provider, err := oidc.NewProvider(ctx, "http://127.0.0.1:8080/realms/dupman")
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC provider: %w", err)
	}

	return &ClientCredentialsProvider{
		config: clientcredentials.Config{
			TokenURL:     provider.Endpoint().TokenURL,
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Scopes:       scopes,
		},
	}, nil
}

func (prov *ClientCredentialsProvider) Retrieve() (*oauth2.Token, error) {
	if prov.IsExpired() {
		ctx := context.Background()

		token, err := prov.config.Token(ctx)
		if err != nil {
			return nil, fmt.Errorf("unable to get token: %w", err)
		}

		prov.token = token
	}

	return prov.token, nil
}

func (prov *ClientCredentialsProvider) IsExpired() bool {
	return !prov.token.Valid()
}
