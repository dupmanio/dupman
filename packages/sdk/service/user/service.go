package user

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// User provides the API operation methods for working with /user routes.
type User struct {
	client *resty.Client
}

// New creates a new instance of the User client with a session.
//
// Example:
//
//	// Create new session with config.
//	sess := session.New(&dupman.Config{...})
//
//	// Create a User client from just a session.
//	svc := user.New(sess)
func New(sess *session.Session) *User {
	return &User{
		client: client.NewUserAPIClient(sess),
	}
}

// Create creates new user.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := user.New(sess)
//
//	// Create new user.
//	account, err := svc.Create(&dto.UserOnCreate{...})
func (svc *User) Create(payload *dto.UserOnCreate) (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	resp, err := svc.client.R().
		SetResult(&response).
		SetBody(payload).
		Post("/user")
	if err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}

// Update updates the user.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := user.New(sess)
//
//	// Update the user.
//	account, err := svc.Update(&dto.UserOnUpdate{...})
func (svc *User) Update(payload *dto.UserOnUpdate) (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	resp, err := svc.client.R().
		SetResult(&response).
		SetBody(payload).
		Patch("/user")
	if err != nil {
		return nil, fmt.Errorf("unable to update User: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}

// GetContactInfo gets user contact info.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := user.New(sess)
//
//	// Get contact info.
//	contactInfo, err := svc.GetContactInfo(...)
func (svc *User) GetContactInfo(id uuid.UUID) (*dto.ContactInfo, error) {
	var response *dto.HTTPResponse[*dto.ContactInfo]

	resp, err := svc.client.R().
		SetResult(&response).
		SetPathParam("id", id.String()).
		Get("/user/contact-info/{id}")
	if err != nil {
		return nil, fmt.Errorf("unable to get User Contact Info: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}

// Me current authenticated user's data.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := user.New(sess)
//
//	// Get current user.
//	contactInfo, err := svc.Me()
func (svc *User) Me() (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	resp, err := svc.client.R().
		SetResult(&response).
		Get("/user/me")
	if err != nil {
		return nil, fmt.Errorf("unable to get User Data: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}
