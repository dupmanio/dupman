package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/inspector"
	"github.com/chromedp/chromedp"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChromeService struct {
	timeout int
	options []chromedp.RemoteAllocatorOption
	config  *config.Config
	logger  *zap.Logger
	ot      *otel.OTel
}

func NewChromeService(conf *config.Config, logger *zap.Logger, ot *otel.OTel) *ChromeService {
	svc := &ChromeService{
		timeout: conf.Chrome.Timeout,
		config:  conf,
		logger:  logger,
		ot:      ot,
		options: []chromedp.RemoteAllocatorOption{},
	}

	return svc
}

func (svc *ChromeService) Screenshot(ctx context.Context, url string, id uuid.UUID) ([]byte, error) {
	var screenshotBuffer []byte

	ctx, span := svc.ot.GetSpanForFunctionCall(ctx, 1)
	defer span.End()

	allocatorCtx, allocatorCancel := chromedp.NewRemoteAllocator(ctx, svc.config.Chrome.RemoteURL, svc.options...)
	defer allocatorCancel()

	browserCtx, cancelBrowserCtx := chromedp.NewContext(allocatorCtx)
	defer cancelBrowserCtx()

	tabCtx, cancelTabCtx := context.WithTimeout(browserCtx, time.Duration(svc.timeout)*time.Second)
	defer cancelTabCtx()

	// Run the initial browser.
	if err := chromedp.Run(browserCtx); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to run browser", err, otel.WebsiteID(id))

		return nil, fmt.Errorf("unable to run browser: %w", err)
	}

	// Prevent browser crashes from locking the context (prevents hanging).
	chromedp.ListenTarget(browserCtx, func(ev interface{}) {
		if _, ok := ev.(*inspector.EventTargetCrashed); ok {
			cancelBrowserCtx()
		}
	})

	chromedp.ListenTarget(tabCtx, func(ev interface{}) {
		if _, ok := ev.(*inspector.EventTargetCrashed); ok {
			cancelTabCtx()
		}
	})

	// Take a screenshot of current tab.
	if err := chromedp.Run(tabCtx, svc.makeScreenshot(url, &screenshotBuffer)); err != nil {
		svc.ot.ErrorEvent(ctx, "Unable to screenshot Website", err, otel.WebsiteID(id))

		return nil, fmt.Errorf("unable to screenshot Website: %w", err)
	}

	return screenshotBuffer, nil
}

func (svc *ChromeService) makeScreenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.EmulateViewport(int64(svc.config.Chrome.ResolutionX), int64(svc.config.Chrome.ResolutionY)),
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(res),
	}
}
