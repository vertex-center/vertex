package types

import (
	"github.com/google/uuid"
)

type EventListener interface {
	OnEvent(e Event)
	GetUUID() uuid.UUID
}
