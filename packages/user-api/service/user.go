package service

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/common/vault"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/user-api/model"
	"github.com/dupmanio/dupman/packages/user-api/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
	vaultSvc *vault.Vault
	ot       *otel.OTel
}

func NewUserService(userRepo *repository.UserRepository, vaultSvc *vault.Vault, ot *otel.OTel) *UserService {
	return &UserService{
		userRepo: userRepo,
		vaultSvc: vaultSvc,
		ot:       ot,
	}
}

func (svc *UserService) Create(ctx context.Context, entity *model.User) (*model.User, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	if user := svc.userRepo.FindByID(ctx, entity.ID.String()); user != nil {
		err := errors.ErrUserAlreadyExists
		svc.ot.ErrorEvent(ctx, "User already exists", err)

		return nil, err
	}

	if err := svc.userRepo.Create(ctx, entity); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to create User", err)

		return nil, fmt.Errorf("unable to create User: %w", err)
	}

	if err := svc.vaultSvc.CreateUserTransitKey(ctx, entity); err != nil {
		return nil, fmt.Errorf("unable to create User Vault key: %w", err)
	}

	return entity, nil
}

func (svc *UserService) Update(ctx context.Context, entity *model.User) (*model.User, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	if user := svc.userRepo.FindByID(ctx, entity.ID.String()); user == nil {
		err := errors.ErrUserDoesNotExist
		svc.ot.ErrorEvent(ctx, "User does not exist", err)

		return nil, err
	}

	if err := svc.userRepo.Update(ctx, entity); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to update User", err)

		return nil, fmt.Errorf("unable to update User: %w", err)
	}

	return entity, nil
}

func (svc *UserService) GetSingle(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	user := svc.userRepo.FindByID(ctx, userID.String())
	if user == nil {
		err := errors.ErrUserDoesNotExist
		svc.ot.ErrorEvent(ctx, "User does not exist", err)

		return nil, err
	}

	return user, nil
}
