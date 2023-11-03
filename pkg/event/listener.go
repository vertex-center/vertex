package event

import (
	"github.com/google/uuid"
)

type Listener interface {
	OnEvent(e interface{})
	GetUUID() uuid.UUID
}
