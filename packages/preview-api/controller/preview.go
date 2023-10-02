package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/preview-api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreviewController struct {
	chromeSvc *service.ChromeService
	httpSvc   *commonServices.HTTPService
	dupmanSvc *service.DupmanAPIService
}

func NewPreviewController(
	chromeSvc *service.ChromeService,
	httpSvc *commonServices.HTTPService,
	dupmanSvc *service.DupmanAPIService,
) (*PreviewController, error) {
	return &PreviewController{
		chromeSvc: chromeSvc,
		httpSvc:   httpSvc,
		dupmanSvc: dupmanSvc,
	}, nil
}

func (ctrl *PreviewController) Preview(ctx *gin.Context) {
	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid website ID: %s", err))

		return
	}

	sess, err := ctrl.dupmanSvc.CreateSession(ctx.GetHeader("Authorization"))
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusUnauthorized, err.Error())

		return
	}

	websiteInstance, err := ctrl.dupmanSvc.GetWebsite(sess, websiteID)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	screen, err := ctrl.chromeSvc.Screenshot(websiteInstance.URL)
	if err != nil {
		ctrl.httpSvc.HTTPError(ctx, http.StatusInternalServerError, err.Error())

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, dto.Preview{
		URL: "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(screen),
	})
}
