package controller

import (
	"errors"
	"fmt"
	"net/http"

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

type WebsiteController struct {
	httpSvc    *commonServices.HTTPService
	websiteSvc *service.WebsiteService
}

func NewWebsiteController(
	httpSvc *commonServices.HTTPService,
	websiteSvc *service.WebsiteService,
) (*WebsiteController, error) {
	return &WebsiteController{httpSvc: httpSvc, websiteSvc: websiteSvc}, nil
}

func (ctrl *WebsiteController) Create(ctx *gin.Context) {
	var (
		entity = &model.Website{}

		payload  *dto.WebsiteOnCreate
		response dto.WebsiteOnCreateResponse
	)

	if err := ctx.ShouldBind(&payload); err != nil {
		ctrl.httpSvc.HTTPValidationError(ctx, err)

		return
	}

	_ = copier.Copy(&entity, &payload)

	website, err := ctrl.websiteSvc.Create(entity, ctx)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusCreated, response)
}

func (ctrl *WebsiteController) GetAll(ctx *gin.Context) {
	var response dto.WebsitesOnResponse

	pager := pagination.Paginate(ctx)

	websites, err := ctrl.websiteSvc.GetAllForCurrentUser(ctx, pager)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	_ = copier.Copy(&response, &websites)

	ctrl.httpSvc.HTTPPaginatedResponse(ctx, http.StatusOK, response, pager)
}

func (ctrl *WebsiteController) GetSingle(ctx *gin.Context) {
	var response dto.WebsiteOnResponse

	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid website ID: %s", err))

		return
	}

	website, err := ctrl.websiteSvc.GetSingleForCurrentUser(ctx, websiteID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, domainErrors.ErrWebsiteNotFound) || errors.Is(err, domainErrors.ErrAccessIsForbidden) {
			statusCode = http.StatusNotFound
			err = domainErrors.ErrWebsiteNotFound
		}

		ctrl.httpSvc.HTTPError(ctx, statusCode, err.Error())

		return
	}

	_ = copier.Copy(&response, &website)

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, response)
}
