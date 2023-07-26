package service

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (svc *UserService) CreateIfNotExists(userID string) (*dto.UserAccount, error) {
	var response dto.UserAccount

	if user := svc.userRepo.FindByID(userID); user != nil {
		_ = copier.Copy(&response, &user)

		return &response, nil
	}

	var newUser model.User

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("unable to parse User ID: %w", err)
	}

	newUser.ID = userUUID
	if err = svc.userRepo.Create(&newUser); err != nil {
		return nil, fmt.Errorf("unable to parse create User: %w", err)
	}

	_ = copier.Copy(&response, &newUser)

	return &response, nil
}

func (svc *UserService) CurrentUserID(ctx *gin.Context) string {
	return ctx.GetString(constant.UserIDKey)
}

func (svc *UserService) CurrentUser(ctx *gin.Context) *model.User {
	return svc.userRepo.FindByID(svc.CurrentUserID(ctx))
}
