package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsiteService struct {
	websiteRepo  *repository.WebsiteRepository
	userSvc      *UserService
	userRepo     *repository.UserRepository
	messengerSvc *MessengerService
	logger       *zap.Logger
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	userSvc *UserService,
	userRepo *repository.UserRepository,
	messengerSvc *MessengerService,
	logger *zap.Logger,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo:  websiteRepo,
		userSvc:      userSvc,
		userRepo:     userRepo,
		messengerSvc: messengerSvc,
		logger:       logger,
	}
}

func (svc *WebsiteService) Create(entity *model.Website, ctx *gin.Context) (*model.Website, error) {
	currentUser := svc.userSvc.CurrentUser(ctx)
	entity.UserID = currentUser.ID

	if err := svc.websiteRepo.Create(entity, currentUser.KeyPair.PublicKey); err != nil {
		return nil, fmt.Errorf("unable to create website: %w", err)
	}

	if err := svc.messengerSvc.SendScanWebsiteMessage(entity); err != nil {
		svc.logger.Error(
			"Unable to publish message to Scanner",
			zap.String("websiteID", entity.ID.String()),
			zap.Error(err),
		)
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

func (svc *WebsiteService) DeleteForCurrentUser(
	ctx *gin.Context,
	websiteID uuid.UUID,
) error {
	currentUser := svc.userSvc.CurrentUser(ctx)
	if err := svc.websiteRepo.DeleteByIDAndUserID(websiteID, currentUser.ID); err != nil {
		return fmt.Errorf("unable to delete websites: %w", err)
	}

	return nil
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

func (svc *WebsiteService) Update(
	website *model.Website,
	updates *dto.WebsiteOnUpdate,
	ctx *gin.Context,
) (*model.Website, error) {
	currentUser := svc.userSvc.CurrentUser(ctx)
	fieldsToUpdate := []string{"URL"}
	website.URL = updates.URL

	if updates.Token != "" {
		website.Token = sqltype.WebsiteToken(updates.Token)

		fieldsToUpdate = append(fieldsToUpdate, "Token")
	}

	if err := svc.websiteRepo.Update(website, fieldsToUpdate, currentUser.KeyPair.PublicKey); err != nil {
		return nil, fmt.Errorf("unable to update website: %w", err)
	}

	return website, nil
}

func (svc *WebsiteService) UpdateStatus(
	website *model.Website,
	newStatus model.Status,
	updates []model.Update,
) (*model.Website, error) {
	if err := svc.websiteRepo.ClearUpdates(website); err != nil {
		return nil, fmt.Errorf("unable to delete Website Updates: %w", err)
	}

	oldStatus := website.Status
	website.Status = newStatus
	website.Status.ID = oldStatus.ID
	website.Status.CreatedAt = oldStatus.CreatedAt

	if newStatus.State == dto.StatusStateNeedsUpdate && updates != nil && len(updates) != 0 {
		website.Updates = updates
	}

	if err := svc.websiteRepo.UpdateStatus(website); err != nil {
		return nil, fmt.Errorf("unable to update Website status: %w", err)
	}

	if err := svc.messengerSvc.SendStatusChangeNotificationMessage(website, oldStatus, newStatus, updates); err != nil {
		svc.logger.Error(
			"Unable to send Status Change notification",
			zap.String("websiteID", website.ID.String()),
			zap.Error(err),
		)
	}

	return website, nil
}
