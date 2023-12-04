package main

import (
	"context"
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/common/otel"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/messenger"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/scheduler"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/version"
)

func process(ctx context.Context) error {
	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("unable to create config: %w", err)
	}

	loggerInst, err := logger.New(conf.Env, conf.AppName, version.Version, conf.LogPath)
	if err != nil {
		return fmt.Errorf("unable to create logger: %w", err)
	}

	ot, err := otel.NewOTel(
		conf.Env,
		conf.AppName,
		version.Version,
		conf.Telemetry.CollectorURL,
		loggerInst,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize Telemetry service: %w", err)
	}

	mess, err := messenger.NewMessengerService(loggerInst, conf, ot)
	if err != nil {
		return fmt.Errorf("unable to create Messenger Service: %w", err)
	}
	defer mess.Close()

	schedulerInst, err := scheduler.New(conf, loggerInst, mess, ot)
	if err != nil {
		return fmt.Errorf("unable to create instance of updater: %w", err)
	}

	if err = schedulerInst.Process(ctx); err != nil {
		return fmt.Errorf("unable to process websites: %w", err)
	}

	return nil
}

func main() {
	ctx := context.Background()
	if err := process(ctx); err != nil {
		panic(err)
	}
}
