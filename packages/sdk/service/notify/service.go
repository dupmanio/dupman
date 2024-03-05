package notify

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/dupmanio/dupman/packages/sdk/service"
)

// Notify provides the API operation methods for working with notify service.
type Notify struct {
	service.Base
}

// New creates a new instance of the Notify service.
func New(conf dupman.Config) *Notify {
	svc := new(Notify)

	svc.SetConfig(conf)
	svc.SetClient(client.NewNotifyClient(conf))

	return svc
}

// Create creates new notification.
func (svc *Notify) Create(
	payload *dto.NotificationOnCreate,
	options ...service.Option,
) (*dto.NotificationOnResponse, error) {
	var response *dto.HTTPResponse[*dto.NotificationOnResponse]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	service.ApplyOptions(req, options)

	resp, err := req.
		SetResult(&response).
		SetBody(payload).
		Post("/notification")
	if err != nil {
		return nil, fmt.Errorf("unable to create Notification: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}
