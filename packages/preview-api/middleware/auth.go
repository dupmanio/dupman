package middleware

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/auth"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authHandler *auth.Handler
	httpSvc     *commonServices.HTTPService
}

func NewAuthMiddleware(
	config *config.Config,
	httpSvc *commonServices.HTTPService,
) (*AuthMiddleware, error) {
	handler, err := auth.NewHandler(auth.HandlerOptions{OauthIssuer: config.OAuth.Issuer})
	if err != nil {
		return nil, fmt.Errorf("unable to initialize auth handler: %w", err)
	}

	return &AuthMiddleware{
		authHandler: handler,
		httpSvc:     httpSvc,
	}, nil
}

func (mid *AuthMiddleware) Setup() {}

func (mid *AuthMiddleware) RequiresAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions {
			ctx.Next()

			return
		}

		claims, err := mid.authHandler.ExtractAuthClaims(ctx.GetHeader("Authorization"))
		if err != nil {
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		mid.authHandler.StoreAuthData(ctx, claims)

		ctx.Next()
	}
}

func (mid *AuthMiddleware) RequiresRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions || mid.authHandler.HasRoles(ctx, roles) {
			ctx.Next()

			return
		}

		mid.httpSvc.HTTPError(ctx, http.StatusForbidden, errors.ErrAccessIsForbidden.Error())
	}
}

func (mid *AuthMiddleware) RequiresScope(scopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Method == http.MethodOptions || mid.authHandler.HasScopes(ctx, scopes) {
			ctx.Next()

			return
		}

		mid.httpSvc.HTTPError(ctx, http.StatusForbidden, errors.ErrAccessIsForbidden.Error())
	}
}
