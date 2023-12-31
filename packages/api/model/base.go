package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (entity *Base) BeforeCreate(_ *gorm.DB) error {
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	return nil
}
