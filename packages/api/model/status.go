package model

import (
	"github.com/google/uuid"
)

type Status struct {
	Base

	WebsiteID uuid.UUID
	State     string `gorm:"type:status_state"`
	Info      string
}
