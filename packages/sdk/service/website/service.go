package website

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/internal/errors"
)

// Website provides the API operation methods for working with /website routes.
type Website struct {
	session *session.Session
}

// New creates a new instance of the Website client with a session.
//
// Example:
//
//	// Create new session with config.
//	sess := session.New(&dupman.Config{...})
//
//	// Create a Website client from just a session.
//	svc := website.New(sess)
func New(sess *session.Session) *Website {
	return &Website{session: sess}
}

// Create creates new website.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := website.New(sess)
//
//	// Create new website.
//	newWebsite, err := svc.Create(&dto.WebsiteOnCreate{...})
func (svc *Website) Create(payload *dto.WebsiteOnCreate) (*dto.WebsiteOnCreateResponse, error) {
	var response *dto.HTTPResponse[*dto.WebsiteOnCreateResponse]

	resp, err := svc.session.Client.R().
		SetResult(&response).
		SetBody(payload).
		Post("/website")
	if err != nil {
		return nil, fmt.Errorf("unable to create Website: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}

// GetAll gets user websites.
// Returns paginated response. You can specify the page argument.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := website.New(sess)
//
//	// Get first page of websites.
//	response, err := svc.GetAll(1)
func (svc *Website) GetAll(page int) (*dto.WebsitesOnResponse, *pagination.Pagination, error) {
	var response *dto.HTTPResponse[*dto.WebsitesOnResponse]

	resp, err := svc.session.Client.R().
		SetResult(&response).
		SetQueryParam("limit", "50").
		SetQueryParam("page", fmt.Sprintf("%d", page)).
		Get("/website")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch Websites: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, nil, errors.NewHTTPError(resp)
	}

	return response.Data, response.Pagination, nil
}
