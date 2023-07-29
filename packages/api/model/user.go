package model

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Base

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
