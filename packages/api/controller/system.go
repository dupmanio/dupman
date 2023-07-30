package controller

import (
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/dto"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SystemController struct {
	httpSvc    *service.HTTPService
	websiteSvc *service.WebsiteService
}

func NewSystemController(
	httpSvc *service.HTTPService,
	websiteSvc *service.WebsiteService,
) (*SystemController, error) {
	return &SystemController{httpSvc: httpSvc, websiteSvc: websiteSvc}, nil
}

func (ctrl *SystemController) GetWebsites(ctx *gin.Context) {
	publicKey := ctx.GetHeader(constant.PublicKeyHeaderKey)
	if publicKey == "" {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("Header '%s' is missing", constant.PublicKeyHeaderKey))

		return
	}

	pager := pagination.Paginate(ctx)

	websites, err := ctrl.websiteSvc.GetAllWithToken(pager, publicKey)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, websites, pager)
}

func (ctrl *SystemController) PutWebsiteUpdates(ctx *gin.Context) {
	var payload dto.Updates

	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid website ID: %s", err))

		return
	}

	if err = ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	updates, code, err := ctrl.websiteSvc.CreateUpdates(websiteID, payload)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, code, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, code, updates)
}
