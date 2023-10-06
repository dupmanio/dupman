package service

import (
	"fmt"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/dupman/session"
	"github.com/dupmanio/dupman/packages/sdk/service/system"
	"github.com/dupmanio/dupman/packages/sdk/service/user"
	"github.com/dupmanio/dupman/packages/sdk/service/website"
)

type DupmanAPIService struct {
	SystemSvc  *system.System
	UserSvc    *user.User
	WebsiteSvc *website.Website
}

func NewDupmanAPIService() *DupmanAPIService {
	return &DupmanAPIService{}
}

func (svc *DupmanAPIService) CreateSession(cred credentials.Provider) error {
	sess, err := session.New(&dupman.Config{Credentials: cred})
	if err != nil {
		return fmt.Errorf("unable to create dupman session: %w", err)
	}

	svc.initializeServices(sess)

	return nil
}

func (svc *DupmanAPIService) initializeServices(sess *session.Session) {
	svc.SystemSvc = system.New(sess)
	svc.UserSvc = user.New(sess)
	svc.WebsiteSvc = website.New(sess)
}
