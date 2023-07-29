package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base

	Username  string `gorm:"unique"`
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Websites  []Website
	KeyPairID uuid.UUID `gorm:"->;<-:create"`
	KeyPair   KeyPair   `gorm:"->;<-:create"`
	Roles     []string  `gorm:"-"`
}

func (entity *User) BeforeCreate(tx *gorm.DB) error {
	if err := entity.Base.BeforeCreate(tx); err != nil {
		return fmt.Errorf("unable to run parent BeforeCreate Callback: %w", err)
	}

	entity.KeyPair = KeyPair{
		PrivateKey: "tmp",
	}

	return nil
}
