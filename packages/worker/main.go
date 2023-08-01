package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/worker/config"
	"github.com/dupmanio/dupman/packages/worker/updater"
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

	upd, err := updater.New(conf, logger)
	if err != nil {
		return fmt.Errorf("unable to create instance of updater: %w", err)
	}

	if err = upd.Process(); err != nil {
		return fmt.Errorf("unable to process websites: %w", err)
	}

	return nil
}

func main() {
	if err := process(); err != nil {
		panic(err)
	}
}
