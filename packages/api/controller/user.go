package controller

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
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
	var payload *dto.UserOnCreate

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	user, err := ctrl.userSvc.CreateIfNotExists(payload)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, user)
}

func (ctrl *UserController) Update(ctx *gin.Context) {
	var payload *dto.UserOnUpdate

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	user, err := ctrl.userSvc.Update(payload)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, user)
}
