package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/domain/dto"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/messenger"
	"github.com/dupmanio/dupman/packages/sdk/dupman"
	"github.com/dupmanio/dupman/packages/sdk/dupman/credentials"
	"github.com/dupmanio/dupman/packages/sdk/service/system"
	"go.uber.org/zap"
)

type Scheduler struct {
	logger        *zap.Logger
	systemService *system.System
	messengerSvc  *messenger.Service
	ot            *otel.OTel
}

func New(
	ctx context.Context,
	conf *config.Config,
	logger *zap.Logger,
	messengerSvc *messenger.Service,
	ot *otel.OTel,
) (*Scheduler, error) {
	cred, err := credentials.NewClientCredentials(ctx, conf.Dupman.ClientID, conf.Dupman.ClientSecret, []string{})
	if err != nil {
		return nil, fmt.Errorf("unable to initiate credentials provider: %w", err)
	}

	dupmanConf := dupman.NewConfig(
		dupman.WithCredentials(cred),
	)

	return &Scheduler{
		logger:        logger,
		systemService: system.New(dupmanConf),
		messengerSvc:  messengerSvc,
		ot:            ot,
	}, nil
}

func (scheduler *Scheduler) Process(ctx context.Context) error {
	scheduler.logger.Info("Starting Scheduler Process")

	currentPage := 1
	totalPages := 1

	var wg sync.WaitGroup

	for currentPage <= totalPages {
		websites, pager, err := scheduler.systemService.GetWebsites(currentPage)
		if err != nil {
			scheduler.logger.Error("unable to get Websites", zap.Error(err))
		}

		totalPages = pager.TotalPages
		scheduler.logger.Info(
			"Fetched single page of websites",
			zap.Int("currentPage", currentPage),
			zap.Int("totalPages", totalPages),
		)
		currentPage++

		for _, website := range *websites {
			wg.Add(1)

			go func(website dto.WebsiteOnSystemResponse) {
				defer wg.Done()

				if err = scheduler.scheduleWebsiteScanning(ctx, website); err != nil {
					scheduler.logger.Error(
						"Unable to schedule website scanning",
						zap.Error(err),
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
					)
				} else {
					scheduler.logger.Info(
						"Website scanning has been scheduled successfully",
						zap.String("websiteID", website.ID.String()),
						zap.String("websiteURL", website.URL),
					)
				}
			}(website)
		}
	}
	wg.Wait()

	return nil
}

func (scheduler *Scheduler) scheduleWebsiteScanning(ctx context.Context, website dto.WebsiteOnSystemResponse) error {
	scheduler.logger.Info(
		"Started processing Website",
		zap.String("websiteID", website.ID.String()),
		zap.String("websiteURL", website.URL),
	)

	if err := scheduler.messengerSvc.SendScanWebsiteMessage(ctx, website); err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	return nil
}
