package service

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
)

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
}

func NewWebsiteService(websiteRepo *repository.WebsiteRepository) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
	}
}

func (svc *WebsiteService) Create(website *model.Website) {
	svc.websiteRepo.Create(website)
}
