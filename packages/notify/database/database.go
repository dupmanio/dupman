package database

import (
	"fmt"
	"net"

	"github.com/dupmanio/dupman/packages/notify/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Database struct {
	*gorm.DB
}

func New(config *config.Config) (*Database, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.Database.User,
		config.Database.Password,
		net.JoinHostPort(config.Database.Host, config.Database.Port),
		config.Database.Database,
	)

	// @todo: implement logging.
	db, err := gorm.Open(postgres.Open(url))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if config.Telemetry.Enabled {
		if err = db.Use(tracing.NewPlugin()); err != nil {
			if err != nil {
				return nil, fmt.Errorf("unable to setup tracing plugin: %w", err)
			}
		}
	}

	return &Database{
		DB: db,
	}, nil
}
