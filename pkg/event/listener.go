package event

import (
	"github.com/google/uuid"
)

type Listener interface {
	OnEvent(e Event) error
	GetUUID() uuid.UUID
}
