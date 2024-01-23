package filter

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/auth"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
)

type RoleFilter struct {
	roles []string
}

func NewRoleFilter(roles ...string) *RoleFilter {
	return &RoleFilter{
		roles: roles,
	}
}

func (filter *RoleFilter) Filter(ctx *gin.Context, authHandler *auth.Handler) (int, error) {
	if !authHandler.HasRoles(ctx, filter.roles) {
		return http.StatusForbidden, domainErrors.ErrAccessIsForbidden
	}

	return http.StatusOK, nil
}
