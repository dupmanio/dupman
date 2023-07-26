package controller

import (
	"net/http"

	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/gin-gonic/gin"
)

type WebsiteController struct {
	httpSvc    *service.HTTPService
	websiteSvc *service.WebsiteService
}

func NewWebsiteController(
	httpSvc *service.HTTPService,
	websiteSvc *service.WebsiteService,
) (*WebsiteController, error) {
	return &WebsiteController{httpSvc: httpSvc, websiteSvc: websiteSvc}, nil
}

func (ctrl *WebsiteController) Create(ctx *gin.Context) {
	var payload *dto.WebsiteOnCreate

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	website, err := ctrl.websiteSvc.Create(payload, ctx)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, website)
}

func (ctrl *WebsiteController) GetAll(ctx *gin.Context) {
	websites, err := ctrl.websiteSvc.GetAll(ctx)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, websites)
}
