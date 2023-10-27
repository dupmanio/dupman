package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationOnCreate struct {
	UserID  uuid.UUID `json:"userID" binding:"required"`
	Type    string    `json:"type" binding:"required"`
	Title   string    `json:"title" binding:"required"`
	Message string    `json:"message" binding:"required"`
}

type NotificationOnResponse struct {
	NotificationOnCreate
	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	Seen      bool      `json:"seen" binding:"required"`
}

type NotificationsOnResponse []NotificationOnResponse
