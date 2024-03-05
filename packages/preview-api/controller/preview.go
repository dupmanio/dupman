package controller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/dupmanio/dupman/packages/common/otel"
	commonServices "github.com/dupmanio/dupman/packages/common/service"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/dupmanio/dupman/packages/preview-api/service"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	sdkService "github.com/dupmanio/dupman/packages/sdk/service"
	"github.com/dupmanio/dupman/packages/sdk/service/website"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PreviewController struct {
	chromeSvc *service.ChromeService
	httpSvc   *commonServices.HTTPService
	config    *config.Config
	ot        *otel.OTel
}

func NewPreviewController(
	chromeSvc *service.ChromeService,
	httpSvc *commonServices.HTTPService,
	config *config.Config,
	ot *otel.OTel,
) (*PreviewController, error) {
	return &PreviewController{
		chromeSvc: chromeSvc,
		httpSvc:   httpSvc,
		config:    config,
		ot:        ot,
	}, nil
}

func (ctrl *PreviewController) Preview(ctx *gin.Context) { //nolint: funlen
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

	websiteSvc := website.New(dupman.NewConfig(
		dupman.WithBaseURL(ctrl.config.ServiceURL.API),
		dupman.WithCredentials(cred),
		dupman.WithOTelEnabled(),
	))

	websiteInstance, err := websiteSvc.Get(websiteID, sdkService.WithContext(ctx))
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
