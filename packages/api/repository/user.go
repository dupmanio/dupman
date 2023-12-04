package repository

import (
	"context"
	"errors"

	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/common/database"
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

func (repo *UserRepository) Create(ctx context.Context, user *model.User) error {
	return repo.db.
		WithContext(ctx).
		Create(user).
		Error
}

func (repo *UserRepository) Update(ctx context.Context, user *model.User) error {
	return repo.db.
		WithContext(ctx).
		Save(user).
		Error
}

func (repo *UserRepository) FindByID(ctx context.Context, id string) *model.User {
	var user model.User

	err := repo.db.
		WithContext(ctx).
		Joins("KeyPair").
		First(&user, "users.id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}
