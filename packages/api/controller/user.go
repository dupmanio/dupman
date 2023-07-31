package controller

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type UserController struct {
	httpSvc *service.HTTPService
	userSvc *service.UserService
}

func NewUserController(httpSvc *service.HTTPService, userSvc *service.UserService) (*UserController, error) {
	return &UserController{httpSvc: httpSvc, userSvc: userSvc}, nil
}

func (ctrl *UserController) Create(ctx *gin.Context) {
	// @todo: refactor: return error if user exists.
	var (
		payload  *dto.UserOnCreate
		response dto.UserAccount
	)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	user, err := ctrl.userSvc.CreateIfNotExists(payload)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &user)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}

func (ctrl *UserController) Update(ctx *gin.Context) {
	var (
		payload  *dto.UserOnUpdate
		response dto.UserAccount
	)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	user, err := ctrl.userSvc.Update(payload)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &user)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}
