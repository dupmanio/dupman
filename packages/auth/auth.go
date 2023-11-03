package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strings"

	"github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	jwt.RegisteredClaims

	Scope       string
	RealmAccess struct {
		Roles []string
	} `json:"realm_access"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (hand *Handler) ExtractAuthClaims(authHeaderValue string) (Claims, error) {
	claims := Claims{}

	accessToken := strings.TrimPrefix(authHeaderValue, "Bearer ")
	if _, _, err := jwt.NewParser().ParseUnverified(accessToken, &claims); err != nil {
		return claims, fmt.Errorf("unable to unmarshal claims : %w", err)
	}

	return claims, nil
}

func (hand *Handler) StoreAuthData(ctx *gin.Context, claims Claims) {
	ctx.Set(constant.UserIDKey, claims.Subject)
	ctx.Set(constant.TokenScopesKey, strings.Split(claims.Scope, " "))
	ctx.Set(constant.TokenRolesKey, claims.RealmAccess.Roles)
}

func (hand *Handler) HasRoles(ctx *gin.Context, roles []string) bool {
	return hand.contextValueContainsArray(ctx, constant.TokenRolesKey, roles)
}

func (hand *Handler) contextValueContainsArray(ctx *gin.Context, valueKey string, needle []string) bool {
	haystack, ok := ctx.Value(valueKey).([]string)
	if !ok {
		haystack = []string{}
	}

	return arrayContainsIntersection(needle, haystack)
}

func (hand *Handler) HasScopes(ctx *gin.Context, scopes []string) bool {
	return hand.contextValueContainsArray(ctx, constant.TokenScopesKey, scopes)
}

func arrayContainsIntersection(arr1, arr2 []string) bool {
	arr1Map := make(map[string]bool)
	for _, str := range arr1 {
		arr1Map[str] = true
	}

	for _, str := range arr2 {
		if arr1Map[str] {
			return true
		}
	}

	return false
}
