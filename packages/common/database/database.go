package database

import (
	"fmt"
	"net"

	"github.com/dupmanio/dupman/packages/common/config"
	logWrapper "github.com/dupmanio/dupman/packages/common/logger/wrapper"
	"github.com/dupmanio/dupman/packages/common/otel"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func New(config config.DatabaseConfig, logger *zap.Logger, ot *otel.OTel) (*Database, error) {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		config.User,
		config.Password,
		net.JoinHostPort(config.Host, config.Port),
		config.Database,
	)
	gormConfig := gorm.Config{
		Logger: logWrapper.NewGormWrapper(logger),
	}

	db, err := gorm.Open(postgres.Open(url), &gormConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = db.Use(ot.GetGormPlugin(config.Database)); err != nil {
		if err != nil {
			return nil, fmt.Errorf("unable to setup tracing plugin: %w", err)
		}
	}

	return &Database{
		DB: db,
	}, nil
}
