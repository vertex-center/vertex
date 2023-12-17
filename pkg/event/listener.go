package event

import (
	"github.com/vertex-center/vertex/common/uuid"
)

type Listener interface {
	OnEvent(e Event) error
	GetUUID() uuid.UUID
}
