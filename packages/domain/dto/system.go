package dto

type WebsiteStatusUpdatePayload struct {
	Status  StatusOnCreate `json:"status" binding:"required"`
	Updates Updates        `json:"updates,omitempty"`
}

type WebsiteStatusUpdateResponse struct {
	Status  StatusOnSystemResponse `json:"status" binding:"required"`
	Updates UpdatesOnResponse      `json:"updates,omitempty"`
}
