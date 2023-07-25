package model

import (
	"github.com/google/uuid"
)

type Website struct {
	Base

	URL    string
	UserID uuid.UUID
}
