package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/dbutils/pagination"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
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
	var response dto.WebsitesOnSystemResponse

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

	_ = copier.Copy(&response, &websites)

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, response, pager)
}

func (ctrl *SystemController) PutWebsiteUpdates(ctx *gin.Context) {
	var (
		entities []model.Update

		payload  dto.Updates
		response dto.UpdatesOnResponse
	)

	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid website ID: %s", err))

		return
	}

	if err = ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	_ = copier.Copy(&entities, &payload)

	updates, err := ctrl.websiteSvc.CreateUpdates(websiteID, entities)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) {
			statusCode = http.StatusNotFound
		}

		ctrl.httpSvc.HTTPError(ctx, statusCode, err.Error())

		return
	}

	_ = copier.Copy(&response, &updates)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}
