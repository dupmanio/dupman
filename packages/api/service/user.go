package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Create(ctx *gin.Context, entity *model.User) (*model.User, error) {
	if user := svc.userRepo.FindByID(ctx.Request.Context(), entity.ID.String()); user != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	if err := svc.userRepo.Create(ctx.Request.Context(), entity); err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	}

	return entity, nil
}

func (svc *UserService) Update(ctx *gin.Context, entity *model.User) (*model.User, error) {
	if user := svc.userRepo.FindByID(ctx.Request.Context(), entity.ID.String()); user == nil {
		return nil, errors.ErrUserDoesNotExist
	}

	if err := svc.userRepo.Update(ctx.Request.Context(), entity); err != nil {
		return nil, fmt.Errorf("unable to update User: %w", err)
	}

	return entity, nil
}

func (svc *UserService) CurrentUser(ctx *gin.Context) *model.User {
	if user, ok := ctx.Get(constant.CurrentUserKey); ok {
		if currentUser, ok := user.(*model.User); ok {
			return currentUser
		}
	}

	return nil
}

func (svc *UserService) GetSingle(
	ctx *gin.Context,
	userID uuid.UUID,
) (*model.User, error) {
	user := svc.userRepo.FindByID(ctx.Request.Context(), userID.String())
	if user == nil {
		return nil, errors.ErrUserDoesNotExist
	}

	return user, nil
}
