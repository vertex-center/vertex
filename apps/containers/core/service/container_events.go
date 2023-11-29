package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
)

func (s *ContainerService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ContainerService) OnEvent(e event.Event) error {
	switch e.(type) {
	case types.EventServerStart:
		s.LoadAll()
	case types.EventServerSetupCompleted:
		go func() {
			s.StartAll()
		}()
	case types.EventServerStop:
		s.StopAll()
	case types.EventServerHardReset:
		s.StopAll()
		s.DeleteAll()
	}
	return nil
}
