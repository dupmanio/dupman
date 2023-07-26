package controller

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	httpSvc *service.HTTPService
	userSvc *service.UserService
}

func NewUserController(httpSvc *service.HTTPService, userSvc *service.UserService) (*UserController, error) {
	return &UserController{httpSvc: httpSvc, userSvc: userSvc}, nil
}

func (ctrl *UserController) Create(ctx *gin.Context) {
	user, err := ctrl.userSvc.CreateIfNotExists(ctx.GetString(constant.UserIDKey))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, user)
}
