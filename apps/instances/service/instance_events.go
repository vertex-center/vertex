package service

import (
	"github.com/google/uuid"
	vtypes "github.com/vertex-center/vertex/types"
)

func (s *InstanceService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *InstanceService) OnEvent(e interface{}) {
	switch e.(type) {
	case vtypes.EventServerStart:
		go func() {
			s.LoadAll()
			s.StartAll()
		}()
	case vtypes.EventServerStop:
		s.StopAll()
	}
}
