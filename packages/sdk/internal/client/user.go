package client

import (
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/go-resty/resty/v2"
)

func NewUserAPIClient(config *dupman.Config) *resty.Client {
	if config.BaseURL == "" {
		config.BaseURL = "http://gateway.dupman.localhost/user-api"
	}

	return getBaseClient(config)
}
