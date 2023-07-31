package session

import (
	"encoding/json"
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
	Config *dupman.Config
	Client *resty.Client
}

// New returns a new Session created from SDK defaults
// or with user provided dupman.Config.
func New(config *dupman.Config) *Session {
	return &Session{
		Config: config,
		Client: createHTTPClient(config),
	}
}

func createHTTPClient(config *dupman.Config) *resty.Client {
	httpClient := resty.New()

	if config.URL == "" {
		config.URL = defaultURL
	}

	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	if config.AccessToken != "" {
		httpClient.SetAuthToken(config.AccessToken)
	}

	httpClient.SetBaseURL(config.URL).
		SetTimeout(config.Timeout).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Type", "application/json").
		SetHeader("User-Agent", "dupman-sdk (https://github.com/dupmanio/dupman/tree/main/packages/sdk)").
		SetDebug(config.Debug).
		SetError(&errors.HTTPError{})

	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal

	return httpClient
}
