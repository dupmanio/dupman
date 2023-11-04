package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/common/pagination"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	domainErrors "github.com/dupmanio/dupman/packages/domain/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type SystemController struct {
	httpSvc    *commonServices.HTTPService
	websiteSvc *service.WebsiteService
}

func NewSystemController(
	httpSvc *commonServices.HTTPService,
	websiteSvc *service.WebsiteService,
) (*SystemController, error) {
	return &SystemController{
		httpSvc:    httpSvc,
		websiteSvc: websiteSvc,
	}, nil
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

func (ctrl *SystemController) UpdateWebsiteStatus(ctx *gin.Context) {
	var (
		statusEntity   model.Status
		updateEntities []model.Update

		payload  dto.WebsiteStatusUpdatePayload
		response dto.WebsiteStatusUpdateResponse
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

	_ = copier.Copy(&statusEntity, &payload.Status)
	_ = copier.Copy(&updateEntities, &payload.Updates)

	website, err := ctrl.websiteSvc.GetSingle(websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) {
			statusCode = http.StatusNotFound
		}

		ctrl.httpSvc.HTTPError(ctx, statusCode, err.Error())

		return
	}

	website, err = ctrl.websiteSvc.UpdateStatus(website, statusEntity, updateEntities)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response.Status, &website.Status)
	_ = copier.Copy(&response.Updates, &website.Updates)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}
