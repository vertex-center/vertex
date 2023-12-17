package types

import (
	"github.com/vertex-center/uuid"
)

type (
	Tags []Tag
	Tag  struct {
		ID     uuid.UUID `json:"id"      db:"id"      example:"12"`
		UserID uuid.UUID `json:"user_id" db:"user_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Name   string    `json:"name"    db:"name"    example:"Vertex SQL"`
	}
)
