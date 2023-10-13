package types

import (
	"github.com/google/uuid"
)

type Listener interface {
	OnEvent(e interface{})
	GetUUID() uuid.UUID
}

type TempListener struct {
	uuid    uuid.UUID
	onEvent func(e interface{})
}

func NewTempListener(onEvent func(e interface{})) TempListener {
	return TempListener{
		uuid:    uuid.New(),
		onEvent: onEvent,
	}
}

func (t TempListener) OnEvent(e interface{}) {
	t.onEvent(e)
}

func (t TempListener) GetUUID() uuid.UUID {
	return t.uuid
}

type (
	EventServerStart struct {
		// PostMigrationCommands are commands that should be executed after the server has started.
		// These are migration commands that cannot be executed before the server has started.
		PostMigrationCommands []interface{}
	}

	EventAppReady struct {
		AppID string
	}

	EventServerStop          struct{}
	EventServerHardReset     struct{}
	EventDependenciesUpdated struct{}
)
