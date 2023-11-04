package service

import (
	"encoding/json"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/broker"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/dupmanio/dupman/packages/common/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebsiteService struct {
	websiteRepo *repository.WebsiteRepository
	userSvc     *UserService
	userRepo    *repository.UserRepository
	broker      *broker.RabbitMQ
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	userSvc *UserService,
	userRepo *repository.UserRepository,
	broker *broker.RabbitMQ,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo: websiteRepo,
		userSvc:     userSvc,
		userRepo:    userRepo,
		broker:      broker,
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

	if err := svc.sendStatusChangeNotification(website, oldStatus, newStatus, updates); err != nil {
		return nil, fmt.Errorf("unable to send Status Change notification: %w", err)
	}

	return website, nil
}

func (svc *WebsiteService) sendStatusChangeNotification(
	website *model.Website,
	oldStatus model.Status,
	newStatus model.Status,
	updates []model.Update,
) error {
	var (
		err                error
		notificationToSend []byte
	)

	if newStatus.State == dto.StatusStateNeedsUpdate && oldStatus.State != dto.StatusStateNeedsUpdate {
		notificationToSend, err = svc.composeNeedsUpdateNotification(website, updates)
		if err != nil {
			return fmt.Errorf("unable to composse Notification: %w", err)
		}
	}

	if newStatus.State == dto.StatusStateScanningFailed && oldStatus.State != dto.StatusStateScanningFailed {
		notificationToSend, err = svc.composeScanningFailedNotification(website, newStatus)
		if err != nil {
			return fmt.Errorf("unable to composse Notification: %w", err)
		}
	}

	if notificationToSend != nil {
		err = svc.broker.PublishToNotify(notificationToSend)
		if err != nil {
			return fmt.Errorf("unable to publish notification: %w", err)
		}
	}

	return nil
}

func (svc *WebsiteService) composeNeedsUpdateNotification(
	website *model.Website,
	updates []model.Update,
) ([]byte, error) {
	updatesMapping := map[string]string{}
	for _, update := range updates {
		updatesMapping[update.Name] = update.RecommendedVersion
	}

	return svc.composeNotification(website.UserID, "WebsiteNeedsUpdates", map[string]any{
		"WebsiteID":  website.ID,
		"WebsiteURL": website.URL,
		"Updates":    updatesMapping,
	})
}

func (svc *WebsiteService) composeScanningFailedNotification(
	website *model.Website,
	status model.Status,
) ([]byte, error) {
	return svc.composeNotification(website.UserID, "WebsiteScanningFailed", map[string]any{
		"WebsiteID":  website.ID,
		"WebsiteURL": website.URL,
		"StatusInfo": status.Info,
	})
}

func (svc *WebsiteService) composeNotification(
	userID uuid.UUID,
	notificationType string,
	notificationMeta map[string]any,
) ([]byte, error) {
	message := dto.NotificationMessage{
		UserID: userID,
		Type:   notificationType,
		Meta:   notificationMeta,
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return jsonMessage, nil
}
