package dto

import (
	"time"

	"github.com/google/uuid"
)

type WebsiteOnCreate struct {
	URL   string `json:"url" binding:"required,url"`
	Token string `json:"token" binding:"required"`
}

type WebsiteOnCreateResponse struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
	URL       string    `json:"url" binding:"required"`
}

type WebsiteOnResponse struct {
	WebsiteOnCreateResponse
	Status  StatusOnWebsitesResponse `json:"status"`
	Updates UpdatesOnResponse        `json:"updates"`
}

type WebsitesOnResponse []WebsiteOnResponse

type WebsiteOnSystemResponse struct {
	ID     uuid.UUID `json:"id" binding:"required"`
	UserID uuid.UUID `json:"userID" binding:"required"`
	URL    string    `json:"url" binding:"required"`
	Token  string    `json:"token" binding:"required"`
}

type WebsitesOnSystemResponse []WebsiteOnSystemResponse

type WebsiteOnUpdate struct {
	URL   string `json:"url" binding:"required,url"`
	Token string `json:"token,omitempty"`
}
