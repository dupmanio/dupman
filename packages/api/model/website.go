package model

import (
	"github.com/google/uuid"
)

type Website struct {
	Base

	URL     string
	Token   string
	UserID  uuid.UUID
	Updates []Update
	Status  Status
}
