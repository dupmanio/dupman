package service

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/pagination"
	commonService "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsiteService struct {
	websiteRepo  *repository.WebsiteRepository
	authSvc      *commonService.AuthService
	userRepo     *repository.UserRepository
	messengerSvc *MessengerService
	logger       *zap.Logger
	ot           *otel.OTel
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	authSvc *commonService.AuthService,
	userRepo *repository.UserRepository,
	messengerSvc *MessengerService,
	logger *zap.Logger,
	ot *otel.OTel,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo:  websiteRepo,
		authSvc:      authSvc,
		userRepo:     userRepo,
		messengerSvc: messengerSvc,
		logger:       logger,
		ot:           ot,
	}
}

func (svc *WebsiteService) Create(ctx context.Context, entity *model.Website) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	currentUser := svc.authSvc.CurrentUser(ctx)
	entity.UserID = currentUser.ID

	if err := svc.websiteRepo.Create(ctx, entity, currentUser.KeyPair.PublicKey); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to create Website", err)

		return nil, fmt.Errorf("unable to create Website: %w", err)
	}

	if err := svc.messengerSvc.SendScanWebsiteMessage(ctx, entity); err != nil {
		svc.ot.LogErrorEvent(ctx, "Unable to publish message to Scanner", err, otel.WebsiteID(entity.ID))
	}

	return entity, nil
}

func (svc *WebsiteService) GetAllForCurrentUser(
	ctx context.Context,
	pagination *pagination.Pagination,
) ([]model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	currentUser := svc.authSvc.CurrentUser(ctx)

	websites, err := svc.websiteRepo.FindByUserID(ctx, currentUser.ID.String(), pagination)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, fmt.Errorf("unable to load Websites: %w", err)
	}

	return websites, nil
}

func (svc *WebsiteService) GetSingleForCurrentUser(
	ctx context.Context,
	websiteID uuid.UUID,
) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	currentUser := svc.authSvc.CurrentUser(ctx)

	website := svc.websiteRepo.FindByID(ctx, websiteID.String())
	if website == nil {
		err := errors.ErrWebsiteNotFound
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, err
	}

	if website.UserID != currentUser.ID {
		err := errors.ErrAccessIsForbidden
		svc.ot.ErrorEvent(ctx, "Access Denied", err)

		return nil, err
	}

	return website, nil
}

func (svc *WebsiteService) DeleteForCurrentUser(
	ctx context.Context,
	websiteID uuid.UUID,
) error {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	currentUser := svc.authSvc.CurrentUser(ctx)
	if err := svc.websiteRepo.DeleteByIDAndUserID(ctx, websiteID, currentUser.ID); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to delete Websites", err)

		return fmt.Errorf("unable to delete Websites: %w", err)
	}

	return nil
}

func (svc *WebsiteService) GetSingle(
	ctx context.Context,
	websiteID uuid.UUID,
) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	website := svc.websiteRepo.FindByID(ctx, websiteID.String())
	if website == nil {
		err := errors.ErrWebsiteNotFound
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, err
	}

	return website, nil
}

func (svc *WebsiteService) GetAllWithToken(
	ctx context.Context,
	pagination *pagination.Pagination,
	publicKey string,
) ([]model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	websites, err := svc.websiteRepo.FindAll(ctx, pagination)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, fmt.Errorf("unable to load Websites: %w", err)
	}

	for i := 0; i < len(websites); i++ {
		// @todo: Implement user key caching.
		user := svc.userRepo.FindByID(ctx, websites[i].UserID.String())

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
	ctx context.Context,
	website *model.Website,
	updates *dto.WebsiteOnUpdate,
) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	currentUser := svc.authSvc.CurrentUser(ctx)
	fieldsToUpdate := []string{"URL"}
	website.URL = updates.URL

	if updates.Token != "" {
		website.Token = sqltype.WebsiteToken(updates.Token)

		fieldsToUpdate = append(fieldsToUpdate, "Token")
	}

	if err := svc.websiteRepo.Update(
		ctx,
		website,
		fieldsToUpdate,
		currentUser.KeyPair.PublicKey,
	); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to update Website", err)

		return nil, fmt.Errorf("unable to update Website: %w", err)
	}

	return website, nil
}

func (svc *WebsiteService) UpdateStatus(
	ctx context.Context,
	website *model.Website,
	newStatus model.Status,
	updates []model.Update,
) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	if err := svc.websiteRepo.ClearUpdates(ctx, website); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to delete Website", err)

		return nil, fmt.Errorf("unable to delete Website Updates: %w", err)
	}

	oldStatus := website.Status
	website.Status = newStatus
	website.Status.ID = oldStatus.ID
	website.Status.CreatedAt = oldStatus.CreatedAt

	if newStatus.State == dto.StatusStateNeedsUpdate && updates != nil && len(updates) != 0 {
		website.Updates = updates
	}

	if err := svc.websiteRepo.UpdateStatus(ctx, website); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to update Website Status", err)

		return nil, fmt.Errorf("unable to update Website Status: %w", err)
	}

	if err := svc.messengerSvc.SendStatusChangeNotificationMessage(
		ctx,
		website,
		oldStatus,
		newStatus,
		updates,
	); err != nil {
		svc.ot.LogErrorEvent(ctx, "Unable to send status change notification", err, otel.WebsiteID(website.ID))
	}

	return website, nil
}
