package service

import (
	"github.com/google/uuid"
	ev "github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/pkg/event"
)

func (s *containerService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *containerService) OnEvent(e event.Event) error {
	switch e.(type) {
	case ev.ServerStart:
		s.LoadAll()
	case ev.ServerSetupCompleted:
		go func() {
			s.StartAll()
		}()
	case ev.ServerStop:
		s.StopAll()
	}
	return nil
}
