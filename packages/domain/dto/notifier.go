package dto

import "github.com/google/uuid"

type NotificationMessage struct {
	UserID uuid.UUID
	Type   string
	Meta   map[string]string
}
