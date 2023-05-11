package repository

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type EventInMemoryRepository struct {
	listeners *map[uuid.UUID]types.Listener
}

func NewEventInMemoryRepository() EventInMemoryRepository {
	return EventInMemoryRepository{
		listeners: &map[uuid.UUID]types.Listener{},
	}
}

func (r *EventInMemoryRepository) AddListener(l types.Listener) {
	(*r.listeners)[l.GetUUID()] = l
}

func (r *EventInMemoryRepository) RemoveListener(l types.Listener) {
	delete(*r.listeners, l.GetUUID())
}

func (r *EventInMemoryRepository) Send(e interface{}) {
	for _, l := range *r.listeners {
		l.OnEvent(e)
	}
}
