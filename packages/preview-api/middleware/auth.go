package middleware

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/auth"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authHandler *auth.Handler
	httpSvc     *commonServices.HTTPService
}

func NewAuthMiddleware(
	httpSvc *commonServices.HTTPService,
) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		authHandler: auth.NewHandler(),
		httpSvc:     httpSvc,
	}, nil
}

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
