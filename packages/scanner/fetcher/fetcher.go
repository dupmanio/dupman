package fetcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/scanner/model"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	timeout    = 30 * time.Second
	retryCount = 3
)

type Fetcher struct {
	client *resty.Client
	logger *zap.Logger
}

func New(logger *zap.Logger) *Fetcher {
	httpClient := resty.New().
		SetTimeout(timeout).
		SetRetryCount(retryCount).
		SetHeader("User-Agent", "dupman-scanner (https://github.com/dupmanio/dupman/tree/main/packages/scanner)")

	return &Fetcher{client: httpClient, logger: logger}
}

func (fetcher *Fetcher) Fetch(url string, id uuid.UUID, token string) ([]model.Update, error) {
	var response *model.Status

	fetcher.logger.Info(
		"Fetching website updates",
		zap.String("websiteID", id.String()),
	)

	resp, err := fetcher.client.R().
		SetResult(&response).
		SetHeader("X-Dupman-Token", token).
		Get(fmt.Sprintf("%s/dupman/status", url))
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Website Updates: %w", err)
	}

	if resp.StatusCode() == http.StatusOK {
		return response.Updates, nil
	}

	if resp.StatusCode() == http.StatusNotFound {
		return nil, errors.ErrNoDupmanEndpoint
	} else if resp.StatusCode() == http.StatusUnauthorized || resp.StatusCode() == http.StatusForbidden {
		return nil, errors.ErrDupmanEndpointAccessDenied
	}

	return nil, errors.ErrUnableToFetchUpdates
}
