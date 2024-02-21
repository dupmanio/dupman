package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/go-resty/resty/v2"
)

type Service interface {
	// SetConfig sets the Service configuration.
	SetConfig(conf *dupman.Config)
	// GetConfig gets the Service configuration.
	GetConfig() *dupman.Config

	// SetClient sets the HTTP Client for Service.
	SetClient(client *resty.Client)
	// GetClient gets Service HTTP Client instance.
	GetClient() *resty.Client

	// Request creates instance of HTTP Request.
	Request() (*resty.Request, error)
}

type Base struct {
	config *dupman.Config
	client *resty.Client
}

func (svc *Base) SetConfig(conf *dupman.Config) {
	svc.config = conf
}

func (svc *Base) GetConfig() *dupman.Config {
	return svc.config
}

func (svc *Base) SetClient(client *resty.Client) {
	svc.client = client
}

func (svc *Base) GetClient() *resty.Client {
	return svc.client
}

func (svc *Base) Request() (*resty.Request, error) {
	if svc.GetConfig().Credentials == nil {
		return svc.GetClient().R(), nil
	}

	token, err := svc.GetConfig().Credentials.Retrieve()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve auth token: %w", err)
	}

	return svc.GetClient().
		SetAuthToken(token.AccessToken).
		R(), nil
}
