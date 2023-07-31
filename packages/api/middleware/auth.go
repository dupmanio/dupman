package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	oauthProvider *oidc.Provider
	httpSvc       *service.HTTPService
	userRepo      *repository.UserRepository
	userSvc       *service.UserService
}

func NewAuthMiddleware(
	config *config.Config,
	httpSvc *service.HTTPService,
	userRepo *repository.UserRepository,
	userSvc *service.UserService,
) (*AuthMiddleware, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.OAuth.Issuer)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize OIDC Provider: %w", err)
	}

	return &AuthMiddleware{
		oauthProvider: provider,
		httpSvc:       httpSvc,
		userRepo:      userRepo,
		userSvc:       userSvc,
	}, nil
}

func (mid *AuthMiddleware) Setup() {}

func (mid *AuthMiddleware) RequiresAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()

			return
		}

		accessToken := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
		if accessToken == "" {
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, errors.ErrAuthorizationRequired.Error())

			return
		}

		verifier := mid.oauthProvider.Verifier(&oidc.Config{SkipClientIDCheck: true})

		token, err := verifier.Verify(ctx, accessToken)
		if err != nil {
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		var claims struct {
			Sub         string
			Scope       string
			RealmAccess struct {
				Roles []string
			} `json:"realm_access"`
		}

		if err = token.Claims(&claims); err != nil {
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		user := mid.userRepo.FindByID(claims.Sub)
		if user == nil {
			user = &model.User{}
			user.ID, _ = uuid.Parse(claims.Sub)
		}

		user.Roles = claims.RealmAccess.Roles

		ctx.Set(constant.CurrentUserKey, user)
		ctx.Set(constant.TokenScopesKey, strings.Split(claims.Scope, " "))

		ctx.Next()
	}
}

func (mid *AuthMiddleware) RequiresRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()

			return
		}

		currentUser := mid.userSvc.CurrentUser(ctx)
		if arrayContainsIntersection(roles, currentUser.Roles) {
			ctx.Next()

			return
		}

		mid.httpSvc.HTTPError(ctx, http.StatusForbidden, errors.ErrAccessIsForbidden.Error())
	}
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

func (mid *AuthMiddleware) RequiresScope(scopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()

			return
		}

		tokenScopes, ok := ctx.Value(constant.TokenScopesKey).([]string)
		if !ok {
			tokenScopes = []string{}
		}

		if arrayContainsIntersection(scopes, tokenScopes) {
			ctx.Next()

			return
		}

		mid.httpSvc.HTTPError(ctx, http.StatusForbidden, errors.ErrMissingScopes.Error())
	}
}
