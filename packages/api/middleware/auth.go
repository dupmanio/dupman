package middleware

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/repository"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/auth"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	authHandler *auth.Handler
	httpSvc     *commonServices.HTTPService
	userRepo    *repository.UserRepository
	userSvc     *service.UserService
}

func NewAuthMiddleware(
	httpSvc *commonServices.HTTPService,
	userRepo *repository.UserRepository,
	userSvc *service.UserService,
) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		authHandler: auth.NewHandler(),
		httpSvc:     httpSvc,
		userRepo:    userRepo,
		userSvc:     userSvc,
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

		user := mid.userRepo.FindByID(claims.Subject)
		if user == nil {
			user = &model.User{}
			user.ID, _ = uuid.Parse(claims.Subject)
		}

		user.Roles = claims.RealmAccess.Roles
		mid.authHandler.StoreAuthData(ctx, claims)
		ctx.Set(constant.CurrentUserKey, user)
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

		mid.httpSvc.HTTPError(ctx, http.StatusForbidden, errors.ErrMissingScopes.Error())
	}
}
