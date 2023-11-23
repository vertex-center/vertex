package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
)

func (s *ContainerService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ContainerService) OnEvent(e event.Event) {
	switch e.(type) {
	case types.EventServerStart:
		go func() {
			s.LoadAll()
			s.StartAll()
			s.ctx.DispatchEvent(types.EventAppReady{
				AppID: "vx-containers",
			})
		}()
	case types.EventServerStop:
		s.StopAll()
	case types.EventServerHardReset:
		s.StopAll()
		s.DeleteAll()
	}
}
