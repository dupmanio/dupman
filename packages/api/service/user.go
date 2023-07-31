package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/domain/dto"
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

func (svc *UserService) CreateIfNotExists(payload *dto.UserOnCreate) (*dto.UserAccount, error) {
	var (
		entity   model.User
		response dto.UserAccount
	)

	if user := svc.userRepo.FindByID(payload.ID.String()); user != nil {
		_ = copier.Copy(&response, &user)

		return &response, nil
	}

	_ = copier.Copy(&entity, &payload)

	if err := svc.userRepo.Create(&entity); err != nil {
		return nil, fmt.Errorf("unable to create User: %w", err)
	}

	_ = copier.Copy(&response, &entity)

	return &response, nil
}

func (svc *UserService) Update(payload *dto.UserOnUpdate) (*dto.UserAccount, error) {
	var (
		entity   model.User
		response dto.UserAccount
	)

	_ = copier.Copy(&entity, &payload)

	if err := svc.userRepo.Update(&entity); err != nil {
		return nil, fmt.Errorf("unable to update User: %w", err)
	}

	_ = copier.Copy(&response, &entity)

	return &response, nil
}

func (svc *UserService) CurrentUser(ctx *gin.Context) *model.User {
	if user, ok := ctx.Get(constant.CurrentUserKey); ok {
		if currentUser, ok := user.(*model.User); ok {
			return currentUser
		}
	}

	return nil
}
