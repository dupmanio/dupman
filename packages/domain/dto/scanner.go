package dto

import "github.com/google/uuid"

type ScanWebsiteMessage struct {
	WebsiteID    uuid.UUID `json:"websiteID"`
	WebsiteURL   string    `json:"websiteURL"`
	WebsiteToken string    `json:"websiteToken"`
}
