package notify

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/go-resty/resty/v2"
)

// Notify provides the API operation methods for working with notify service.
type Notify struct {
	client *resty.Client
}

// New creates a new instance of the Notify client with a session.
//
// Example:
//
//	// Create new session with config.
//	sess := session.New(&dupman.Config{...})
//
//	// Create a Notify client from just a session.
//	svc := notify.New(sess)
func New(sess *session.Session) *Notify {
	return &Notify{
		client: client.NewNotifyClient(sess),
	}
}

// Create creates new notification.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := notify.New(sess)
//
//	// Create new notification.
//	notification, err := svc.Create(&dto.NotificationOnCreate{...})
func (svc *Notify) Create(payload *dto.NotificationOnCreate) (*dto.NotificationOnResponse, error) {
	var response *dto.HTTPResponse[*dto.NotificationOnResponse]

	resp, err := svc.client.R().
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
