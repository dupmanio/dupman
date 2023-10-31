package client

import (
	"encoding/json"
	"time"

	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/internal/errors"
	"github.com/go-resty/resty/v2"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultRetryCount = 3
)

func getBaseClient(config *dupman.Config, accessToken string) *resty.Client {
	if config.Timeout == 0 {
		config.Timeout = defaultTimeout
	}

	if config.RetryCount == 0 {
		config.RetryCount = defaultRetryCount
	}

	httpClient := resty.New()
	httpClient.
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetDebug(config.Debug).
		SetAuthToken(accessToken).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Type", "application/json").
		SetHeader("User-Agent", "dupman-sdk (https://github.com/dupmanio/dupman/tree/main/packages/sdk)").
		SetError(&errors.HTTPError{})

	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal

	return httpClient
}
