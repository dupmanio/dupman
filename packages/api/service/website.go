package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
	userSvc     *UserService
}

func NewWebsiteService(websiteRepo *repository.WebsiteRepository, userSvc *UserService) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
		userSvc:     userSvc,
	}
}

func (svc *WebsiteService) Create(payload *dto.WebsiteOnCreate, ctx *gin.Context) (*dto.WebsiteOnResponse, error) {
	var (
		entity   model.Website
		response dto.WebsiteOnResponse
	)

	_ = copier.Copy(&entity, &payload)

	currentUser := svc.userSvc.CurrentUser(ctx)
	entity.UserID = currentUser.ID

	if err := svc.websiteRepo.Create(&entity, currentUser.KeyPair.PublicKey); err != nil {
		return nil, fmt.Errorf("unable to create website: %w", err)
	}

	_ = copier.Copy(&response, &entity)

	return &response, nil
}

func (svc *WebsiteService) GetAll(
	ctx *gin.Context,
	pagination *pagination.Pagination,
) (*dto.WebsitesOnResponse, error) {
	response := dto.WebsitesOnResponse{}

	websites, err := svc.websiteRepo.FindByUserID(svc.userSvc.CurrentUserID(ctx), pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to get websites: %w", err)
	}

	_ = copier.Copy(&response, &websites)

	return &response, nil
}
