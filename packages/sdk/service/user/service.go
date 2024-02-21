package user

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/dupmanio/dupman/packages/sdk/service"
	"github.com/google/uuid"
)

// User provides the API operation methods for working with user service.
type User struct {
	service.Base
}

// New creates a new instance of the User service.
func New(conf *dupman.Config) *User {
	svc := new(User)

	svc.SetConfig(conf)
	svc.SetClient(client.NewUserAPIClient(conf))

	return svc
}

// Create creates new user.
func (svc *User) Create(payload *dto.UserOnCreate) (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
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
func (svc *User) Update(payload *dto.UserOnUpdate) (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
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
func (svc *User) GetContactInfo(id uuid.UUID) (*dto.ContactInfo, error) {
	var response *dto.HTTPResponse[*dto.ContactInfo]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
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

// Me gets current authenticated user's data.
func (svc *User) Me() (*dto.UserAccount, error) {
	var response *dto.HTTPResponse[*dto.UserAccount]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
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
