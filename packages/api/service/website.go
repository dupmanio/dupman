package service

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/common/ory/keto"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/pagination"
	commonService "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/common/vault"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WebsiteService struct {
	websiteRepo  *repository.WebsiteRepository
	authSvc      *commonService.AuthService
	messengerSvc *MessengerService
	logger       *zap.Logger
	ot           *otel.OTel
	vault        *vault.Vault
	keto         *keto.Keto
}

func NewWebsiteService(
	websiteRepo *repository.WebsiteRepository,
	authSvc *commonService.AuthService,
	messengerSvc *MessengerService,
	logger *zap.Logger,
	ot *otel.OTel,
	vault *vault.Vault,
	keto *keto.Keto,
) *WebsiteService {
	return &WebsiteService{
		websiteRepo:  websiteRepo,
		authSvc:      authSvc,
		messengerSvc: messengerSvc,
		logger:       logger,
		ot:           ot,
		vault:        vault,
		keto:         keto,
	}
}

func (svc *WebsiteService) Create(ctx context.Context, entity *model.Website) (*model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	entity.UserID = svc.authSvc.CurrentUserID(ctx)

	token, err := svc.vault.EncryptWithUserTransitKey(ctx, entity.UserID, entity.Token)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to Encrypt Website Token", err)

		return nil, fmt.Errorf("unable to Encrypt Website Token: %w", err)
	}

	entity.Token = token

	if err = svc.websiteRepo.Create(ctx, entity); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to create Website", err)

		return nil, fmt.Errorf("unable to create Website: %w", err)
	}

	if err = svc.keto.CreateRelationship(
		ctx,
		"Website",
		entity.ID.String(),
		"owners",
		entity.UserID.String(),
	); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to Create Website Owner Relationship", err)

		return nil, fmt.Errorf("unable to Create Website Owner Relationship: %w", err)
	}

	if err = svc.messengerSvc.SendScanWebsiteMessage(ctx, entity); err != nil {
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

	currentUserID := svc.authSvc.CurrentUserID(ctx)

	websites, err := svc.websiteRepo.FindByUserID(ctx, currentUserID.String(), pagination)
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

	currentUserID := svc.authSvc.CurrentUserID(ctx)

	website := svc.websiteRepo.FindByID(ctx, websiteID.String())
	if website == nil {
		err := errors.ErrWebsiteNotFound
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, err
	}

	if website.UserID != currentUserID {
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

	currentUserID := svc.authSvc.CurrentUserID(ctx)
	if err := svc.websiteRepo.DeleteByIDAndUserID(ctx, websiteID, currentUserID); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to delete Websites", err)

		return fmt.Errorf("unable to delete Websites: %w", err)
	}

	if err := svc.keto.DeleteRelationship(
		ctx,
		"Website",
		websiteID.String(),
		"owners",
		currentUserID.String(),
	); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to Delete Website Owner Relationship", err)

		return fmt.Errorf("unable to Delete Website Owner Relationship: %w", err)
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
) ([]model.Website, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	websites, err := svc.websiteRepo.FindAll(ctx, pagination)
	if err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to load Websites", err)

		return nil, fmt.Errorf("unable to load Websites: %w", err)
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

	currentUserID := svc.authSvc.CurrentUserID(ctx)
	fieldsToUpdate := []string{"URL"}
	website.URL = updates.URL

	if updates.Token != "" {
		token, err := svc.vault.EncryptWithUserTransitKey(ctx, currentUserID, updates.Token)
		if err != nil {
			svc.ot.ErrorEvent(ctx, "Unable to Encrypt Website Token", err)

			return nil, fmt.Errorf("unable to Encrypt Website Token: %w", err)
		}

		website.Token = token

		fieldsToUpdate = append(fieldsToUpdate, "Token")
	}

	if err := svc.websiteRepo.Update(
		ctx,
		website,
		fieldsToUpdate,
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
