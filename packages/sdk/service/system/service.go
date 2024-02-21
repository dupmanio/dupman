package system

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

// System provides the API operation methods for working with /system routes.
type System struct {
	service.Base
}

// New creates a new instance of the System service.
func New(conf *dupman.Config) *System {
	svc := new(System)

	svc.SetConfig(conf)
	svc.SetClient(client.NewAPIClient(conf))

	return svc
}

// GetWebsites gets all websites.
// Returns paginated response. You can specify the page argument.
func (svc *System) GetWebsites(
	page int,
) (*dto.WebsitesOnSystemResponse, *pagination.Pagination, error) {
	var response *dto.HTTPResponse[*dto.WebsitesOnSystemResponse]

	req, err := svc.Request()
	if err != nil {
		return nil, nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
		SetResult(&response).
		SetQueryParam("limit", "50").
		SetQueryParam("page", strconv.Itoa(page)).
		Get("/system/websites")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to fetch Websites: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, nil, errors.NewHTTPError(resp)
	}

	return response.Data, response.Pagination, nil
}

// UpdateWebsiteStatus updates website status and creates update entities.
func (svc *System) UpdateWebsiteStatus(
	websiteID uuid.UUID,
	status *dto.Status,
	updates *dto.Updates,
) (*dto.WebsiteStatusUpdateResponse, error) {
	var (
		response *dto.HTTPResponse[*dto.WebsiteStatusUpdateResponse]
		payload  = dto.WebsiteStatusUpdatePayload{
			Status: *status,
		}
	)

	if status.State == dto.StatusStateNeedsUpdate {
		payload.Updates = *updates
	}

	req, err := svc.Request()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize request: %w", err)
	}

	resp, err := req.
		SetResult(&response).
		SetBody(payload).
		SetPathParam("websiteId", websiteID.String()).
		Post("/system/websites/{websiteId}/status")
	if err != nil {
		return nil, fmt.Errorf("unable to create Website Updates: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}
