package model

import (
	"github.com/google/uuid"
)

type Notification struct {
	Base

	UserID  uuid.UUID
	Type    string
	Title   string
	Message string
	Seen    bool
}
