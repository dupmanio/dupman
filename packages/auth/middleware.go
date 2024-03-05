package auth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	sdkErrors "github.com/dupmanio/dupman/packages/sdk/errors"
	sdkService "github.com/dupmanio/dupman/packages/sdk/service"
	"github.com/dupmanio/dupman/packages/sdk/service/user"
	"github.com/gin-gonic/gin"
)

type Filter interface {
	Filter(ctx *gin.Context, handler *Handler) (httpCode int, err error)
}

type HTTPErrorHandlerFunc func(ctx *gin.Context, httpCode int, err any)

type Middleware struct {
	handler *Handler

	httpErrorHandler HTTPErrorHandlerFunc
	fetchUserData    bool
	filters          []Filter
}

func NewMiddleware(options ...Option) *Middleware {
	mid := &Middleware{
		handler: NewHandler(),

		filters: make([]Filter, 0),
	}

	mid.applyOptions(options)

	return mid
}

func (mid *Middleware) applyOptions(options []Option) {
	for _, opt := range options {
		opt(mid)
	}
}

func (mid *Middleware) Handler(options ...Option) gin.HandlerFunc {
	instance := *mid
	instance.applyOptions(options)

	return instance.getGinHandler()
}

func (mid *Middleware) getGinHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			userData *dto.UserAccount
			code     int
		)

		claims, rawToken, err := mid.checkAuthHeader(ctx)
		if err != nil {
			mid.httpErrorHandler(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		if mid.fetchUserData {
			if userData, code, err = mid.getUserData(ctx, rawToken); err != nil {
				mid.httpErrorHandler(ctx, code, err.Error())

				return
			}
		}

		mid.handler.StoreAuthData(ctx, claims, userData)
		// @todo: run additional store callbacks.

		if code, err := mid.applyFilters(ctx); err != nil {
			mid.httpErrorHandler(ctx, code, err.Error())

			return
		}

		ctx.Next()
	}
}

func (mid *Middleware) checkAuthHeader(ctx *gin.Context) (Claims, string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return Claims{}, "", domainErrors.ErrAuthorizationRequired
	}

	claims, err := mid.handler.ExtractAuthClaims(authHeader)
	if err != nil {
		return Claims{}, "", err
	}

	return claims, authHeader, nil
}

func (mid *Middleware) getUserData(ctx *gin.Context, token string) (*dto.UserAccount, int, error) {
	svc, err := mid.createDupmanUserSvc(token)
	if err != nil {
		return nil, http.StatusInternalServerError, domainErrors.ErrSomethingWentWrong
	}

	data, err := svc.Me(sdkService.WithContext(ctx))
	if err != nil {
		var sdkErr *sdkErrors.HTTPError
		if errors.As(err, &sdkErr) {
			if sdkErr.Code == http.StatusUnauthorized || sdkErr.Code == http.StatusForbidden {
				return nil, sdkErr.Code, err //nolint: wrapcheck
			}
		}

		return nil, http.StatusInternalServerError, domainErrors.ErrSomethingWentWrong
	}

	return data, http.StatusOK, nil
}

func (mid *Middleware) createDupmanUserSvc(token string) (*user.User, error) {
	cred, err := credentials.NewRawTokenCredentials(token)
	if err != nil {
		return nil, fmt.Errorf("unable to create Dupman Credentials: %w", err)
	}

	// @todo: pass User API URL.
	return user.New(dupman.NewConfig(
		dupman.WithCredentials(cred),
		dupman.WithOTelEnabled(),
	)), nil
}

func (mid *Middleware) applyFilters(ctx *gin.Context) (int, error) {
	for _, filterInstance := range mid.filters {
		if code, err := filterInstance.Filter(ctx, mid.handler); err != nil {
			return code, err //nolint: wrapcheck
		}
	}

	return http.StatusOK, nil
}
