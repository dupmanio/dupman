package auth

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/auth/constant"
	"github.com/gin-gonic/gin"
)

func SetUserID(ctx *gin.Context) {
	userID := ctx.GetHeader(constant.UserIDHeader)

	if userID == "" {
		// @todo: change response.
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})

		return
	}

	ctx.Set(constant.UserIDKey, userID)

	ctx.Next()
}
