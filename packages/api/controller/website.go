package controller

import (
	"errors"
	"fmt"
	"net/http"

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

type WebsiteController struct {
	httpSvc    *commonServices.HTTPService
	websiteSvc *service.WebsiteService
	ot         *otel.OTel
}

func NewWebsiteController(
	httpSvc *commonServices.HTTPService,
	websiteSvc *service.WebsiteService,
	ot *otel.OTel,
) (*WebsiteController, error) {
	return &WebsiteController{
		httpSvc:    httpSvc,
		websiteSvc: websiteSvc,
		ot:         ot,
	}, nil
}

func (ctrl *WebsiteController) Create(ctx *gin.Context) {
	var (
		entity = &model.Website{}

		payload  *dto.WebsiteOnCreate
		response dto.WebsiteOnCreateResponse
	)

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationErrorWithOTelLog(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	website, err := ctrl.websiteSvc.Create(ctx, entity)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to create Website", http.StatusInternalServerError, err)

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.ot.LogInfoEvent(ctx, "Website has been created successfully", otel.WebsiteID(website.ID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}

func (ctrl *WebsiteController) GetAll(ctx *gin.Context) {
	var response dto.WebsitesOnResponse

	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)
	pager := pagination.Paginate(ctx)

	websites, err := ctrl.websiteSvc.GetAllForCurrentUser(ctx, pager)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Websites", http.StatusInternalServerError, err)

		return
	}

	_ = copier.Copy(&response, &websites)

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, response, pager)
}

func (ctrl *WebsiteController) GetSingle(ctx *gin.Context) {
	var response dto.WebsiteOnResponse

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

	website, err := ctrl.websiteSvc.GetSingleForCurrentUser(ctx, websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) || errors.Is(err, domainErrors.ErrAccessIsForbidden) {
			statusCode = http.StatusNotFound
			err = domainErrors.ErrWebsiteNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Website", statusCode, err, otel.WebsiteID(websiteID))

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}

func (ctrl *WebsiteController) Update(ctx *gin.Context) {
	var (
		payload  *dto.WebsiteOnUpdate
		response dto.WebsiteOnCreateResponse
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

	website, err := ctrl.websiteSvc.GetSingleForCurrentUser(ctx, websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) || errors.Is(err, domainErrors.ErrAccessIsForbidden) {
			statusCode = http.StatusNotFound
			err = domainErrors.ErrWebsiteNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Website", statusCode, err, otel.WebsiteID(websiteID))

		return
	}

	if website, err = ctrl.websiteSvc.Update(ctx, website, payload); err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to update Website",
			http.StatusInternalServerError,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.ot.LogInfoEvent(ctx, "Website has been updated successfully", otel.WebsiteID(website.ID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}

func (ctrl *WebsiteController) Delete(ctx *gin.Context) {
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

	_, err = ctrl.websiteSvc.GetSingleForCurrentUser(ctx, websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) || errors.Is(err, domainErrors.ErrAccessIsForbidden) {
			statusCode = http.StatusNotFound
			err = domainErrors.ErrWebsiteNotFound
		}

		ctrl.httpSvc.HTTPErrorWithOTelLog(ctx, "Unable to load Website", statusCode, err, otel.WebsiteID(websiteID))

		return
	}

	err = ctrl.websiteSvc.DeleteForCurrentUser(ctx, websiteID)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to delete Website",
			http.StatusInternalServerError,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	ctrl.ot.LogInfoEvent(ctx, "Website has been deleted successfully", otel.WebsiteID(websiteID))
	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, nil)
}
