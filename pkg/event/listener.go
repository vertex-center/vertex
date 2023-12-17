package event

import (
	"github.com/vertex-center/uuid"
)

type Listener interface {
	OnEvent(e Event) error
	GetUUID() uuid.UUID
}
