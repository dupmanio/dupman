package dto

import (
	"time"

	"github.com/google/uuid"
)

type WebsiteOnCreate struct {
	URL   string `json:"url" binding:"required,url"`
	Token string `json:"token" binding:"required"`
}

type WebsiteOnResponse struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
	URL       string    `json:"url" binding:"required"`
}

type WebsitesOnResponse []WebsiteOnResponse