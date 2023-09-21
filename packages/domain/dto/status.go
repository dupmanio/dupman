package dto

import (
	"time"

	"github.com/google/uuid"
)

type Status struct {
	State string `json:"state" binding:"required"`
	Info  string `json:"info"`
}

type StatusOnSystemResponse struct {
	Status

	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}
