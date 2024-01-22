package system

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/internal/client"
	"github.com/dupmanio/dupman/packages/sdk/internal/errors"
	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

// System provides the API operation methods for working with /system routes.
type System struct {
	client *resty.Client
}

// New creates a new instance of the System client with a session.
//
// Example:
//
//	// Create new session with config.
//	sess := session.New(&dupman.Config{...})
//
//	// Create a System client from just a session.
//	svc := system.New(sess)
func New(sess *session.Session) *System {
	return &System{
		client: client.NewAPIClient(sess),
	}
}

// GetWebsites gets all websites.
// Returns paginated response. You can specify the page argument.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := system.New(sess)
//
//	// Get first page of websites.
//	websites, pager, err := svc.GetWebsites(1)
func (svc *System) GetWebsites(
	page int,
) (*dto.WebsitesOnSystemResponse, *pagination.Pagination, error) {
	var response *dto.HTTPResponse[*dto.WebsitesOnSystemResponse]

	resp, err := svc.client.R().
		SetResult(&response).
		SetQueryParam("limit", "50").
		SetQueryParam("page", fmt.Sprintf("%d", page)).
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
//
// Example:
//
//	// Create new instance of service using session.
//	svc := system.New(sess)
//
//	// Create updates.
//	status, err := svc.UpdateWebsiteStatus(websiteID, &dto.Status{...}, &dto.Updates{...})
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

	resp, err := svc.client.R().
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
