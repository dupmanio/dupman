package dto

import "github.com/google/uuid"

type NotificationMessage struct {
	UserID uuid.UUID      `json:"userID"`
	Type   string         `json:"type"`
	Meta   map[string]any `json:"meta"`
}
