package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/service/website"
	"github.com/google/uuid"
)

type DupmanAPIService struct{}

func NewDupmanAPIService() *DupmanAPIService {
	return &DupmanAPIService{}
}

func (svc *DupmanAPIService) CreateSession(accessToken string) (*session.Session, error) {
	cred, err := credentials.NewRawTokenCredentials(accessToken)
	if err != nil {
		return nil, fmt.Errorf("unable to initiate credentials provider: %w", err)
	}

	sess, err := session.New(&dupman.Config{Credentials: cred})
	if err != nil {
		return nil, fmt.Errorf("unable to create dupman sesssion: %w", err)
	}

	return sess, nil
}

func (svc *DupmanAPIService) GetWebsite(sess *session.Session, websiteID uuid.UUID) (*dto.WebsiteOnResponse, error) {
	websiteSvc := website.New(sess)

	websiteInstance, err := websiteSvc.Get(websiteID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch website: %w", err)
	}

	return websiteInstance, nil
}
