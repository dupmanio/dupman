package model

import (
	sqltype "github.com/dupmanio/dupman/packages/api/sql/type"
	"github.com/google/uuid"
)

type Website struct {
	Base

	URL     string
	Token   sqltype.WebsiteToken
	UserID  uuid.UUID
	Updates []Update
	Status  Status
}
