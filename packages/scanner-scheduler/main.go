package main

import (
	"fmt"

	"github.com/dupmanio/dupman/packages/scanner-scheduler/broker"
	"github.com/dupmanio/dupman/packages/scanner-scheduler/config"
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

	brk, err := broker.NewRabbitMQ(conf, logger)
	if err != nil {
		return fmt.Errorf("unable to create RabbitMQ Broker: %w", err)
	}

	schedulerInst, err := scheduler.New(conf, logger, brk)
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
