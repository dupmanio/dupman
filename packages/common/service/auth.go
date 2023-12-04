package service

import (
	"context"

	// @todo: refactor: remove API service dependency.
	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	authConstant "github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/google/uuid"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (svc *AuthService) CurrentUserID(ctx context.Context) uuid.UUID {
	var userID uuid.UUID
	if userIDRaw, ok := ctx.Value(authConstant.UserIDKey).(string); ok {
		userID, _ = uuid.Parse(userIDRaw)
	}

	return userID
}

func (svc *AuthService) CurrentUser(ctx context.Context) *model.User {
	if user, ok := ctx.Value(constant.CurrentUserKey).(*model.User); ok {
		return user
	}

	return nil
}
