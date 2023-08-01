package fetcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/worker/model"
	"github.com/go-resty/resty/v2"
)

const (
	timeout    = 30 * time.Second
	retryCount = 3
)

type Fetcher struct {
	client *resty.Client
}

func New() *Fetcher {
	httpClient := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetHeader("User-Agent", "dupman-worker (https://github.com/dupmanio/dupman/tree/main/packages/worker)")

	return &Fetcher{client: httpClient}
}

func (fetcher *Fetcher) Fetch(url, token string) ([]model.Update, error) {
	var response *model.Status

	resp, err := fetcher.client.R().
		SetResult(&response).
		SetHeader("X-Dupman-Token", token).
		Get(fmt.Sprintf("%s/dupman/status", url))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Website Updates: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.ErrUnableToFetchUpdates
	}

	return response.Updates, nil
}
