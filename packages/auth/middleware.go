package auth

import (
	"net/http"

	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
)

type Filter interface {
	Filter(ctx *gin.Context, handler *Handler) (httpCode int, err error)
}

type HTTPErrorHandlerFunc func(ctx *gin.Context, httpCode int, err any)

type Middleware struct {
	handler *Handler

	httpErrorHandler HTTPErrorHandlerFunc
	callUserService  bool
	filters          []Filter
}

func NewMiddleware(options ...Option) *Middleware {
	mid := &Middleware{
		handler: NewHandler(),

		callUserService: true,
		filters:         make([]Filter, 0),
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
		claims, err := mid.checkAuthHeader(ctx)
		if err != nil {
			mid.httpErrorHandler(ctx, http.StatusUnauthorized, err.Error())

			return
		}

		// @todo: optionally call User API and get current user's data.
		// if mid.callUserService {}

		mid.handler.StoreAuthData(ctx, claims)
		// @todo: run additional store callbacks.

		if code, err := mid.applyFilters(ctx); err != nil {
			mid.httpErrorHandler(ctx, code, err.Error())

			return
		}

		ctx.Next()
	}
}

func (mid *Middleware) checkAuthHeader(ctx *gin.Context) (Claims, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return Claims{}, domainErrors.ErrAuthorizationRequired
	}

	claims, err := mid.handler.ExtractAuthClaims(authHeader)
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}

func (mid *Middleware) applyFilters(ctx *gin.Context) (int, error) {
	for _, filterInstance := range mid.filters {
		if code, err := filterInstance.Filter(ctx, mid.handler); err != nil {
			return code, err //nolint: wrapcheck
		}
	}

	return http.StatusOK, nil
}
