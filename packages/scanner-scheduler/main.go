package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/common/logger"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/messenger"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/scheduler"
)

func process() error {
	conf, err := config.New()
	if err != nil {
		return fmt.Errorf("unable to create config: %w", err)
	}

	loggerInst, err := logger.New(conf.Env, conf.AppName, "1.0.0", conf.LogPath)
	if err != nil {
		return fmt.Errorf("unable to create logger: %w", err)
	}

	mess, err := messenger.NewMessengerService(loggerInst, conf)
	if err != nil {
		return fmt.Errorf("unable to create Messenger Service: %w", err)
	}
	defer mess.Close()

	schedulerInst, err := scheduler.New(conf, loggerInst, mess)
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
