package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/messenger"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/scheduler"
	"go.uber.org/zap"
)

func process() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("unable to create logger: %w", err)
	}

	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("unable to create config: %w", err)
	}

	mess, err := messenger.NewMessengerService(logger, conf)
	if err != nil {
		return fmt.Errorf("unable to create Messenger Service: %w", err)
	}
	defer mess.Close()

	schedulerInst, err := scheduler.New(conf, logger, mess)
	if err != nil {
		return fmt.Errorf("unable to create instance of updater: %w", err)
	}

	if err = schedulerInst.Process(); err != nil {
		return fmt.Errorf("unable to process websites: %w", err)
	}

	return nil
}

func main() {
	if err := process(); err != nil {
		panic(err)
	}
}
