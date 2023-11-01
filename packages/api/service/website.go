package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
	userSvc     *UserService
	userRepo    *repository.UserRepository
	updateRepo  *repository.UpdateRepository
	statusRepo  *repository.StatusRepository
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	userSvc *UserService,
	userRepo *repository.UserRepository,
	updateRepo *repository.UpdateRepository,
	statusRepo *repository.StatusRepository,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
		userSvc:     userSvc,
		userRepo:    userRepo,
		updateRepo:  updateRepo,
		statusRepo:  statusRepo,
	}
}

func (svc *WebsiteService) Create(entity *model.Website, ctx *gin.Context) (*model.Website, error) {
	currentUser := svc.userSvc.CurrentUser(ctx)
	entity.UserID = currentUser.ID

	if err := svc.websiteRepo.Create(entity, currentUser.KeyPair.PublicKey); err != nil {
		return nil, fmt.Errorf("unable to create website: %w", err)
	}

	return entity, nil
}

func (svc *WebsiteService) GetAllForCurrentUser(
	ctx *gin.Context,
	pagination *pagination.Pagination,
) ([]model.Website, error) {
	currentUser := svc.userSvc.CurrentUser(ctx)

	websites, err := svc.websiteRepo.FindByUserID(currentUser.ID.String(), pagination)
	if err != nil {
		return nil, fmt.Errorf("unable to get websites: %w", err)
	}

	return websites, nil
}

func (svc *WebsiteService) GetSingleForCurrentUser(
	ctx *gin.Context,
	websiteID uuid.UUID,
) (*model.Website, error) {
	currentUser := svc.userSvc.CurrentUser(ctx)

	website := svc.websiteRepo.FindByID(websiteID.String())
	if website == nil {
		return nil, errors.ErrWebsiteNotFound
	}

	if website.UserID != currentUser.ID {
		return nil, errors.ErrAccessIsForbidden
	}

	return website, nil
}

func (svc *WebsiteService) GetSingle(
	websiteID uuid.UUID,
) (*model.Website, error) {
	website := svc.websiteRepo.FindByID(websiteID.String())
	if website == nil {
		return nil, errors.ErrWebsiteNotFound
	}

	return website, nil
}

func (svc *WebsiteService) GetAllWithToken(
	pagination *pagination.Pagination,
	publicKey string,
) ([]model.Website, error) {
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

	return websites, nil
}

func (svc *WebsiteService) CreateUpdates(websiteID uuid.UUID, updates []model.Update) ([]model.Update, error) {
	if website := svc.websiteRepo.FindByID(websiteID.String()); website == nil {
		return nil, errors.ErrWebsiteNotFound
	}

	// @todo: refactor: do not delete and rewrite updates.
	if err := svc.updateRepo.DeleteByWebsiteID(websiteID.String()); err != nil {
		return nil, fmt.Errorf("unable to delete Website Updates: %w", err)
	}

	for i := range updates {
		updates[i].WebsiteID = websiteID

		if err := svc.updateRepo.Create(&updates[i]); err != nil {
			return nil, fmt.Errorf("unable to create Website Update: %w", err)
		}
	}

	return updates, nil
}

func (svc *WebsiteService) UpdateStatus(websiteID uuid.UUID, status *model.Status) (*model.Status, error) {
	if website := svc.websiteRepo.FindByID(websiteID.String()); website == nil {
		return nil, errors.ErrWebsiteNotFound
	}

	if statusEntity := svc.statusRepo.FindByWebsiteID(websiteID.String()); statusEntity != nil {
		statusEntity.State = status.State
		statusEntity.Info = status.Info

		if err := svc.statusRepo.Update(statusEntity); err != nil {
			return nil, fmt.Errorf("unable to update website status: %w", err)
		}

		return statusEntity, nil
	}

	status.WebsiteID = websiteID
	if err := svc.statusRepo.Create(status); err != nil {
		return nil, fmt.Errorf("unable to create website status: %w", err)
	}

	return status, nil
}
