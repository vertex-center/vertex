package services

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
)

type InstanceLogsService struct {
	uuid uuid.UUID

	logsAdapter  types.InstanceLogsAdapterPort
	eventAdapter types.EventAdapterPort
}

func NewInstanceLogsService(logsAdapter types.InstanceLogsAdapterPort, eventAdapter types.EventAdapterPort) InstanceLogsService {
	s := InstanceLogsService{
		uuid: uuid.New(),

		logsAdapter:  logsAdapter,
		eventAdapter: eventAdapter,
	}

	s.eventAdapter.AddListener(&s)

	return s
}

func (s *InstanceLogsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *InstanceLogsService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.logsAdapter.LoadBuffer(uuid)
}

func (s *InstanceLogsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceLog:
		s.logsAdapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	case types.EventInstanceLoaded:
		err := s.logsAdapter.Open(e.InstanceUuid)
		if err != nil {
			log.Error(err)
			return
		}
	}
}
