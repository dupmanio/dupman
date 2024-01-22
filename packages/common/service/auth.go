package service

import (
	"context"

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
