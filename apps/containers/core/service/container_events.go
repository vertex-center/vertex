package service

import (
	"context"

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
		s.LoadAll(context.Background())
	case ev.ServerSetupCompleted:
		go func() {
			s.StartAll(context.Background())
		}()
	case ev.ServerStop:
		s.StopAll(context.Background())
	}
	return nil
}
