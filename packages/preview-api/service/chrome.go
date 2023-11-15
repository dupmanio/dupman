package service

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/inspector"
	"github.com/chromedp/chromedp"
	"github.com/dupmanio/dupman/packages/preview-api/config"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChromeService struct {
	timeout int
	options []chromedp.ExecAllocatorOption
	logger  *zap.Logger
}

func NewChromeService(conf *config.Config, logger *zap.Logger) *ChromeService {
	svc := &ChromeService{
		timeout: conf.Chrome.Timeout,
		logger:  logger,
		options: []chromedp.ExecAllocatorOption{
			chromedp.WindowSize(conf.Chrome.ResolutionX, conf.Chrome.ResolutionY),
			chromedp.DisableGPU,
		},
	}

	svc.options = append(svc.options, chromedp.DefaultExecAllocatorOptions[:]...)

	return svc
}

func (svc *ChromeService) Screenshot(url string, id uuid.UUID) ([]byte, error) {
	var screenshotBuffer []byte

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(context.Background(), svc.options...)
	defer allocatorCancel()

	browserCtx, cancelBrowserCtx := chromedp.NewContext(allocatorCtx)
	defer cancelBrowserCtx()

	tabCtx, cancelTabCtx := context.WithTimeout(browserCtx, time.Duration(svc.timeout)*time.Second)
	defer cancelTabCtx()

	// Run the initial browser.
	if err := chromedp.Run(browserCtx); err != nil {
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
		svc.logger.Error(
			"unable to screenshot Website",
			zap.String("websiteID", id.String()),
			zap.Error(err),
		)

		return nil, fmt.Errorf("unable to screenshot Website: %w", err)
	}

	return screenshotBuffer, nil
}

func (svc *ChromeService) makeScreenshot(url string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.CaptureScreenshot(res),
	}
}
