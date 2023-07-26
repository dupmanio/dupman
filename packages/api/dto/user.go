package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserAccount struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}
