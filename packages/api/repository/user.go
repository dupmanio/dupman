package repository

import (
	"errors"

	"github.com/dupmanio/dupman/packages/api/database"
	"github.com/dupmanio/dupman/packages/api/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *database.Database
	logger *zap.Logger
}

func NewUserRepository(
	db *database.Database,
	logger *zap.Logger,
) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (repo *UserRepository) Setup() {
	repo.logger.Debug("Setting up User repository")

	if err := repo.db.AutoMigrate(&model.User{}); err != nil {
		repo.logger.Error(err.Error())
	}
}

func (repo *UserRepository) Create(user *model.User) error {
	return repo.db.Create(user).Error
}

func (repo *UserRepository) Update(user *model.User) error {
	return repo.db.Save(user).Error
}

func (repo *UserRepository) FindByID(id string) *model.User {
	var user model.User

	err := repo.db.Joins("KeyPair").First(&user, "users.id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}
