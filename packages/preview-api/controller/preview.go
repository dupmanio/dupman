package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/preview-api/service"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreviewController struct {
	chromeSvc *service.ChromeService
	httpSvc   *commonServices.HTTPService
	dupmanSvc *commonServices.DupmanAPIService
	ot        *otel.OTel
}

func NewPreviewController(
	chromeSvc *service.ChromeService,
	httpSvc *commonServices.HTTPService,
	dupmanSvc *commonServices.DupmanAPIService,
	ot *otel.OTel,
) (*PreviewController, error) {
	return &PreviewController{
		chromeSvc: chromeSvc,
		httpSvc:   httpSvc,
		dupmanSvc: dupmanSvc,
		ot:        ot,
	}, nil
}

func (ctrl *PreviewController) Preview(ctx *gin.Context) {
	ctrl.httpSvc.EnrichSpanWithControllerAttributes(ctx)

	websiteID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Invalid Website ID",
			http.StatusBadRequest,
			fmt.Errorf("invalid website ID: %w", err),
			otel.WebsiteID(websiteID),
		)

		return
	}

	cred, err := credentials.NewRawTokenCredentials(ctx.GetHeader("Authorization"))
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to create dupman credentials",
			http.StatusUnauthorized,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	err = ctrl.dupmanSvc.CreateSession(cred)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to create dupman session",
			http.StatusUnauthorized,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	// @todo: add tracing data to sdk request headers.
	websiteInstance, err := ctrl.dupmanSvc.WebsiteSvc.Get(websiteID)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to load Website",
			http.StatusInternalServerError,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	screen, err := ctrl.chromeSvc.Screenshot(ctx, websiteInstance.URL, websiteInstance.ID)
	if err != nil {
		ctrl.httpSvc.HTTPErrorWithOTelLog(
			ctx,
			"Unable to screenshot Website",
			http.StatusInternalServerError,
			err,
			otel.WebsiteID(websiteID),
		)

		return
	}

	ctrl.httpSvc.HTTPResponse(ctx, http.StatusOK, dto.Preview{
		URL: "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(screen),
	})
}
