package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dupmanio/dupman/packages/api/config"
	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	oauthProvider *oidc.Provider
	httpSvc       *service.HTTPService
}

func NewAuthMiddleware(config *config.Config, httpSvc *service.HTTPService) (*AuthMiddleware, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, config.OAuth.Issuer)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize OIDC Provider: %w", err)
	}

	return &AuthMiddleware{
		oauthProvider: provider,
		httpSvc:       httpSvc,
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
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, "authentication required")

			return
		}

		// @todo: refactor!
		verifier := mid.oauthProvider.Verifier(&oidc.Config{SkipClientIDCheck: true})

		token, err := verifier.Verify(ctx, accessToken)
		if err != nil {
			mid.httpSvc.HTTPError(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		claims := map[string]interface{}{}
		_ = token.Claims(&claims)

		ctx.Set(constant.UserIDKey, claims["sub"])

		ctx.Next()
	}
}
