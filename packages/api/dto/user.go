package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserOnCreate struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	FirstName string    `json:"firstName" binding:"required"`
	LastName  string    `json:"lastName" binding:"required"`
}

type UserAccount struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	Username  string    `json:"username" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	FirstName string    `json:"firstName" binding:"required"`
	LastName  string    `json:"lastName" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}
