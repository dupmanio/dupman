package website

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/errors"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/dupmanio/dupman/packages/sdk/service"
	"github.com/google/uuid"
)

// Website provides the API operation methods for working with /website routes.
type Website struct {
	service.Base
}

// New creates a new instance of the Website client with a session.
func New(conf *dupman.Config) *Website {
	svc := new(Website)

	svc.SetConfig(conf)
	svc.SetClient(client.NewAPIClient(conf))

	return svc
}

// Create creates new website.
func (svc *Website) Create(payload *dto.WebsiteOnCreate) (*dto.WebsiteOnCreateResponse, error) {
	var response *dto.HTTPResponse[*dto.WebsiteOnCreateResponse]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
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
func (svc *Website) GetAll(page int) (*dto.WebsitesOnResponse, *pagination.Pagination, error) {
	var response *dto.HTTPResponse[*dto.WebsitesOnResponse]

	req, err := svc.Request()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
		SetResult(&response).
		SetQueryParam("limit", "50").
		SetQueryParam("page", strconv.Itoa(page)).
		Get("/website")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch Websites: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, nil, errors.NewHTTPError(resp)
	}

	return response.Data, response.Pagination, nil
}

// Get gets single website.
func (svc *Website) Get(id uuid.UUID) (*dto.WebsiteOnResponse, error) {
	var response *dto.HTTPResponse[*dto.WebsiteOnResponse]

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
		SetResult(&response).
		SetPathParam("id", id.String()).
		Get("/website/{id}")
	if err != nil {
		return nil, fmt.Errorf("unable to fetch Website: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}
