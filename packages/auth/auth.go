package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Sub         string
	Scope       string
	RealmAccess struct {
		Roles []string
	} `json:"realm_access"`
}

type HandlerOptions struct {
	OauthIssuer string
}

type Handler struct {
	oauthProvider *oidc.Provider
}

func NewHandler(options HandlerOptions) (*Handler, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, options.OauthIssuer)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize OIDC Provider: %w", err)
	}

	return &Handler{
		oauthProvider: provider,
	}, nil
}

func (hand *Handler) ExtractAuthClaims(authHeaderValue string) (Claims, error) {
	var claims Claims

	accessToken := strings.TrimPrefix(authHeaderValue, "Bearer ")
	if accessToken == "" {
		return claims, errors.ErrAuthorizationRequired
	}

	token, err := hand.extractIDToken(accessToken)
	if err != nil {
		return claims, fmt.Errorf("invalid token: %w", err)
	}

	if err = token.Claims(&claims); err != nil {
		return claims, fmt.Errorf("unable to unmarshal claims : %w", err)
	}

	return claims, nil
}

func (hand *Handler) extractIDToken(rawToken string) (*oidc.IDToken, error) {
	ctx := context.Background()
	verifierOptions := &oidc.Config{
		SkipClientIDCheck: true,
	}
	verifier := hand.oauthProvider.Verifier(verifierOptions)

	token, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("unable to verify token: %w", err)
	}

	return token, nil
}

func (hand *Handler) StoreAuthData(ctx *gin.Context, claims Claims) {
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
