package database

import (
	"fmt"
	"net"

	"github.com/dupmanio/dupman/packages/api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	return &Database{
		DB: db,
	}, nil
}
