package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

var errWebsiteNotFound = errors.New("website not found")

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
	userSvc     *UserService
	userRepo    *repository.UserRepository
	updateRepo  *repository.UpdateRepository
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	userSvc *UserService,
	userRepo *repository.UserRepository,
	updateRepo *repository.UpdateRepository,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
		userSvc:     userSvc,
		userRepo:    userRepo,
		updateRepo:  updateRepo,
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

func (svc *WebsiteService) GetAllForCurrentUser(
	ctx *gin.Context,
	pagination *pagination.Pagination,
) (*dto.WebsitesOnResponse, error) {
	response := dto.WebsitesOnResponse{}
	currentUser := svc.userSvc.CurrentUser(ctx)

	websites, err := svc.websiteRepo.FindByUserID(currentUser.ID.String(), pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to get websites: %w", err)
	}

	_ = copier.Copy(&response, &websites)

	return &response, nil
}

func (svc *WebsiteService) GetAllWithToken(
	pagination *pagination.Pagination,
	publicKey string,
) (*dto.WebsitesOnSystemResponse, error) {
	response := dto.WebsitesOnSystemResponse{}

	websites, err := svc.websiteRepo.FindAll(pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to get websites: %w", err)
	}

	for i := 0; i < len(websites); i++ {
		// @todo: Implement user key caching.
		user := svc.userRepo.FindByID(websites[i].UserID.String())

		if rawToken, err := websites[i].Token.Decrypt(user.KeyPair.PrivateKey); err == nil {
			websites[i].Token = sqltype.WebsiteToken(rawToken)

			if tokenEncrypted, err := websites[i].Token.Encrypt(publicKey); err == nil {
				websites[i].Token = sqltype.WebsiteToken(tokenEncrypted)
			} else {
				websites[i].Token = ""
			}
		}
	}

	_ = copier.Copy(&response, &websites)

	return &response, nil
}

func (svc *WebsiteService) CreateUpdates(
	websiteID uuid.UUID,
	payload dto.Updates,
) (*dto.UpdatesOnResponse, int, error) {
	response := dto.UpdatesOnResponse{}

	if website := svc.websiteRepo.FindByID(websiteID.String()); website == nil {
		return nil, http.StatusNotFound, errWebsiteNotFound
	}

	if err := svc.updateRepo.DeleteByWebsiteID(websiteID.String()); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("unable to delete Website Updates: %w", err)
	}

	for i := range payload {
		var (
			entity           = model.Update{}
			updateOnResponse = dto.UpdateOnResponse{}
		)

		_ = copier.Copy(&entity, &payload[i])
		entity.WebsiteID = websiteID

		if err := svc.updateRepo.Create(&entity); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("unable to create Website Update: %w", err)
		}

		_ = copier.Copy(&updateOnResponse, &entity)
		response = append(response, updateOnResponse)
	}

	return &response, http.StatusCreated, nil
}
