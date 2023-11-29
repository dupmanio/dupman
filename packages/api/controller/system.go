package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/api/constant"
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/dupmanio/dupman/packages/common/otel"
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
	ot         *otel.OTel
}

func NewSystemController(
	httpSvc *commonServices.HTTPService,
	websiteSvc *service.WebsiteService,
	ot *otel.OTel,
) (*SystemController, error) {
	return &SystemController{
		httpSvc:    httpSvc,
		websiteSvc: websiteSvc,
		ot:         ot,
	}, nil
}

func (ctrl *SystemController) GetWebsites(ctx *gin.Context) {
	var response dto.WebsitesOnSystemResponse

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	publicKey := ctx.GetHeader(constant.PublicKeyHeaderKey)
	if publicKey == "" {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			fmt.Sprintf("Header '%s' is missing", constant.PublicKeyHeaderKey),
			http.StatusInternalServerError,
			domainErrors.ErrHeaderIsMissing,
		)

		return
	}

	pager := pagination.Paginate(ctx)

	websites, err := ctrl.websiteSvc.GetAllWithToken(ctx, pager, publicKey)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Websites", http.StatusInternalServerError, err)

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

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Invalid Website ID",
			http.StatusBadRequest,
			fmt.Errorf("invalid Website ID: %w", err),
		)

		return
	}

	if err = ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationErrorWithOTelLog(ctx, err)

		return
	}

	_ = copier.Copy(&statusEntity, &payload.Status)
	_ = copier.Copy(&updateEntities, &payload.Updates)

	website, err := ctrl.websiteSvc.GetSingle(ctx, websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) {
			statusCode = http.StatusNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Website", statusCode, err, otel.WebsiteID(websiteID))

		return
	}

	website, err = ctrl.websiteSvc.UpdateStatus(ctx, website, statusEntity, updateEntities)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to update Website Status",
			http.StatusInternalServerError,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	// @todo: refactor.
	_ = copier.Copy(&response.Status, &website.Status)
	_ = copier.Copy(&response.Updates, &website.Updates)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}
