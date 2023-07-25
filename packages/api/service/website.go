package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/jinzhu/copier"
)

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
}

func NewWebsiteService(websiteRepo *repository.WebsiteRepository) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
	}
}

func (svc *WebsiteService) Create(payload *dto.WebsiteOnCreate) (*dto.WebsiteOnResponse, error) {
	var (
		entity   model.Website
		response dto.WebsiteOnResponse
	)

	_ = copier.Copy(&entity, &payload)

	if err := svc.websiteRepo.Create(&entity); err != nil {
		return nil, fmt.Errorf("unable to create website: %w", err)
	}

	_ = copier.Copy(&response, &entity)

	return &response, nil
}

func (svc *WebsiteService) GetAll() (*dto.WebsitesOnResponse, error) {
	var response dto.WebsitesOnResponse

	websites, err := svc.websiteRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("unable to get websites: %w", err)
	}

	_ = copier.Copy(&response, &websites)

	return &response, nil
}
