package credentials

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"golang.org/x/oauth2"
)

// RawTokenProvider implements Credentials Provider
// with raw JWT access token.
type RawTokenProvider struct {
	token *oauth2.Token
}

// NewRawTokenCredentials creates a new instance of the RawTokenProvider.
//
// Example:
//
//	// Create new instance of Credential Provider.
//	cred, err := credentials.NewRawTokenCredentials("eyJhb...j4IIA")
//
//	// Create new session.
//	sess, err := session.New(&dupman.Config{Credentials: cred})
func NewRawTokenCredentials(rawToken string) (*RawTokenProvider, error) {
	ctx := context.Background()
	// @todo: update url!
	provider, err := oidc.NewProvider(ctx, "http://id.dupman.localhost/realms/dupman")
	if err != nil {
		return nil, fmt.Errorf("unable to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
	token := strings.Fields(rawToken)

	idToken, err := verifier.Verify(ctx, token[1])
	if err != nil {
		return nil, fmt.Errorf("unable to verify token: %w", err)
	}

	return &RawTokenProvider{
		token: &oauth2.Token{
			TokenType:   token[0],
			AccessToken: token[1],
			Expiry:      idToken.Expiry,
		},
	}, nil
}

func (prov *RawTokenProvider) Retrieve() (*oauth2.Token, error) {
	if prov.IsExpired() {
		return nil, errors.ErrTokenIsExpired
	}

	return prov.token, nil
}

func (prov *RawTokenProvider) IsExpired() bool {
	return !prov.token.Valid()
}
