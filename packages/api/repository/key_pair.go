package repository

import (
	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"go.uber.org/zap"
)

type KeyPairRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewKeyPairRepository(
	db *database.Database,
	logger *zap.Logger,
) *KeyPairRepository {
	return &KeyPairRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *KeyPairRepository) Setup() {
	repo.logger.Debug("Setting up KeyPair repository")

	if err := repo.db.AutoMigrate(&model.KeyPair{}); err != nil {
		repo.logger.Error(err.Error())
	}
}
