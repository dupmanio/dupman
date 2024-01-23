package filter

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/auth"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
)

type ScopeFilter struct {
	scopes []string
}

func NewScopeFilter(scopes ...string) *ScopeFilter {
	return &ScopeFilter{
		scopes: scopes,
	}
}

func (filter *ScopeFilter) Filter(ctx *gin.Context, authHandler *auth.Handler) (int, error) {
	if !authHandler.HasScopes(ctx, filter.scopes) {
		return http.StatusForbidden, domainErrors.ErrAccessIsForbidden
	}

	return http.StatusOK, nil
}
