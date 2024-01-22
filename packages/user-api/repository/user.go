package repository

import (
	"context"
	"errors"

	"github.com/dupmanio/dupman/packages/common/database"
	"github.com/dupmanio/dupman/packages/user-api/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{
		db: db,
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
		First(&user, "users.id = ?", id).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return &user
}
