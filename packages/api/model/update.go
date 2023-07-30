package model

import "github.com/google/uuid"

type Update struct {
	Base

	WebsiteID          uuid.UUID
	Name               string
	Title              string
	Link               string
	Type               string
	CurrentVersion     string
	LatestVersion      string
	RecommendedVersion string
	InstallType        string
	Status             int
}
