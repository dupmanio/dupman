package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (svc *UserService) Create(payload *dto.UserOnCreate) (*model.User, error) {
	var entity model.User

	if user := svc.userRepo.FindByID(payload.ID.String()); user != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	_ = copier.Copy(&entity, &payload)

	if err := svc.userRepo.Create(&entity); err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	}

	return &entity, nil
}

func (svc *UserService) Update(payload *dto.UserOnUpdate) (*model.User, error) {
	var entity model.User

	_ = copier.Copy(&entity, &payload)

	if err := svc.userRepo.Update(&entity); err != nil {
		return nil, fmt.Errorf("unable to update User: %w", err)
	}

	return &entity, nil
}

func (svc *UserService) CurrentUser(ctx *gin.Context) *model.User {
	if user, ok := ctx.Get(constant.CurrentUserKey); ok {
		if currentUser, ok := user.(*model.User); ok {
			return currentUser
		}
	}

	return nil
}
