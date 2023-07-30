package dto

import (
	"time"

	"github.com/google/uuid"
)

type Update struct {
	Name               string `json:"name" binding:"required"`
	Title              string `json:"title" binding:"required"`
	Link               string `json:"link" binding:"required"`
	Type               string `json:"type" binding:"required"`
	CurrentVersion     string `json:"currentVersion" binding:"required"`
	LatestVersion      string `json:"latestVersion" binding:"required"`
	RecommendedVersion string `json:"recommendedVersion" binding:"required"`
	InstallType        string `json:"installType" binding:"required"`
	Status             int    `json:"status" binding:"required"`
}

type Updates []Update

type UpdateOnResponse struct {
	Update
	ID        uuid.UUID `json:"id" binding:"required"`
	CreatedAt time.Time `json:"createdAt" binding:"required"`
	UpdatedAt time.Time `json:"updatedAt" binding:"required"`
}

type UpdatesOnResponse []UpdateOnResponse
