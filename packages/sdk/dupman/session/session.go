package session

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"golang.org/x/oauth2"
)

type Session struct {
	Config *dupman.Config
	Token  *oauth2.Token
}

// New returns a new Session created from SDK defaults
// or with user provided dupman.Config.
//
// Example:
//
//	// Create new instance of Credential Provider, e.g ClientCredentials.
//	cred, err := credentials.NewClientCredentials("...", "...", []string{...})
//
//	// Create new session.
//	sess, err := session.New(&dupman.Config{Credentials: cred})
func New(config *dupman.Config) (*Session, error) {
	token, err := config.Credentials.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve credentials: %w", err)
	}

	return &Session{
		Config: config,
		Token:  token,
	}, nil
}
