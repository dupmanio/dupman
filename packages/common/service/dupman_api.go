package service

import (
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/service/notify"
	"github.com/dupmanio/dupman/packages/sdk/service/system"
	"github.com/dupmanio/dupman/packages/sdk/service/user"
	"github.com/dupmanio/dupman/packages/sdk/service/website"
)

type DupmanAPIService struct {
	NotifySvc  *notify.Notify
	SystemSvc  *system.System
	UserSvc    *user.User
	WebsiteSvc *website.Website
}

func NewDupmanAPIService() *DupmanAPIService {
	return &DupmanAPIService{}
}

func (svc *DupmanAPIService) CreateSession(cred credentials.Provider, additionalOptions ...dupman.Option) error {
	conf := dupman.NewConfig(
		append(
			additionalOptions,
			dupman.WithCredentials(cred),
			dupman.WithDebug(true),
		)...,
	)

	svc.initializeServices(conf)

	return nil
}

func (svc *DupmanAPIService) initializeServices(conf dupman.Config) {
	svc.NotifySvc = notify.New(conf)
	svc.SystemSvc = system.New(conf)
	svc.UserSvc = user.New(conf)
	svc.WebsiteSvc = website.New(conf)
}
