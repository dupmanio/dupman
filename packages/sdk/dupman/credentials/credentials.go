package credentials

import (
	"golang.org/x/oauth2"
)

type Provider interface {
	// Retrieve retrieves token from  the source.
	Retrieve() (*oauth2.Token, error)

	// IsExpired checks if token is expired.
	IsExpired() bool
}
