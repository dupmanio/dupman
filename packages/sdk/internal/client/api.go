package client

import (
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/go-resty/resty/v2"
)

func NewAPIClient(config *dupman.Config) *resty.Client {
	if config.BaseURL == "" {
		config.BaseURL = "http://gateway.dupman.localhost/api"
	}

	return getBaseClient(config)
}
