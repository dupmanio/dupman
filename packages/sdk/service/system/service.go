package system

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/internal/errors"
	"github.com/google/uuid"
)

// System provides the API operation methods for working with /system routes.
type System struct {
	session *session.Session
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
	return &System{session: sess}
}

// GetWebsites gets all websites.
// Requires publicKey generated by encryptor.Encryptor.
// Returns paginated response. You can specify the page argument.
//
// Example:
//
//	// Generate RSA Key Pair.
//	enc := encryptor.NewRSAEncryptor()
//	err := enc.GenerateKeyPair()
//	publicKey, err := enc.PublicKey()
//
//	// Create new instance of service using session.
//	svc := system.New(sess)
//
//	// Get first page of websites.
//	websites, pager, err := svc.GetWebsites(publicKey, 1)
func (svc *System) GetWebsites(
	publicKey string,
	page int,
) (*dto.WebsitesOnSystemResponse, *pagination.Pagination, error) {
	var response *dto.HTTPResponse[*dto.WebsitesOnSystemResponse]

	resp, err := svc.session.Client.R().
		SetHeader("X-Public-Key", publicKey).
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

// CreateWebsiteUpdates creates website update entity.
//
// Example:
//
//	// Create new instance of service using session.
//	svc := system.New(sess)
//
//	// Create updates.
//	updates, err := svc.CreateWebsiteUpdates(websiteID, &dto.Updates{...})
func (svc *System) CreateWebsiteUpdates(websiteID uuid.UUID, payload *dto.Updates) (*dto.UpdatesOnResponse, error) {
	var response *dto.HTTPResponse[*dto.UpdatesOnResponse]

	resp, err := svc.session.Client.R().
		SetResult(&response).
		SetBody(payload).
		SetPathParam("websiteId", websiteID.String()).
		Put("/system/websites/{websiteId}/updates")
	if err != nil {
		return nil, fmt.Errorf("unable to create Website Updates: %w", err)
	}

	if resp.StatusCode() != http.StatusCreated {
		return nil, errors.NewHTTPError(resp)
	}

	return response.Data, nil
}
