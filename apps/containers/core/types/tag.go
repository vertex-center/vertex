package types

import "github.com/google/uuid"

type (
	TagID struct{ uuid.UUID }
	Tags  []Tag
	Tag   struct {
		ID     TagID  `json:"id"      db:"id"      example:"12"`
		UserID uint   `json:"user_id" db:"user_id" example:"d1fb743c-f937-4f3d-95b9-1a8475464591"`
		Name   string `json:"name"    db:"name"    example:"Vertex SQL"`
	}
)

func NewTagID() TagID        { return TagID{uuid.New()} }
func (TagID) Type() string   { return "string" }
func (TagID) Format() string { return "uuid" }
