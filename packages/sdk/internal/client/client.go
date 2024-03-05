package client

import (
	"encoding/json"

	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func getBaseClient(config dupman.Config) *resty.Client {
	httpClient := resty.New()
	httpClient.
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetDebug(config.Debug).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Type", "application/json").
		SetHeader("User-Agent", "dupman-sdk (https://github.com/dupmanio/dupman/tree/main/packages/sdk)").
		SetError(&errors.HTTPError{})

	if config.OTelEnabled {
		httpClient.SetTransport(otelhttp.NewTransport(httpClient.GetClient().Transport))
	}

	httpClient.JSONMarshal = json.Marshal
	httpClient.JSONUnmarshal = json.Unmarshal

	return httpClient
}
