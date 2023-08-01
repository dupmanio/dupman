package session

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/internal/errors"
	"github.com/go-resty/resty/v2"
)

const (
	defaultTimeout = 30 * time.Second
	// @todo: update.
	defaultURL = "http://127.0.0.1:8000"
)

type Session struct {
	config *dupman.Config
	Client *resty.Client
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
		config: config,
		Client: createHTTPClient(config, token.AccessToken),
	}, nil
}

func createHTTPClient(config *dupman.Config, accessToken string) *resty.Client {
	httpClient := resty.New()

	if config.URL == "" {
		config.URL = defaultURL
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	httpClient.SetBaseURL(config.URL).
		SetTimeout(config.Timeout).
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Type", "application/json").
		SetHeader("User-Agent", "dupman-sdk (https://github.com/dupmanio/dupman/tree/main/packages/sdk)").
		SetDebug(config.Debug).
		SetError(&errors.HTTPError{})

	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal

	return httpClient
}
