package services

import (
	"errors"

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
		s.onLogReceived(e)
	case types.EventInstanceLoaded:
		err := s.logsAdapter.Register(e.Instance.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types.EventInstanceDeleted:
		err := s.logsAdapter.Unregister(e.InstanceUUID)
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *InstanceLogsService) onLogReceived(e types.EventInstanceLog) {
	switch e.Kind {
	case types.LogKindDownload:
		var downloads *types.LogLineMessageDownloads
		download := e.Message.(*types.LogLineMessageDownload)

		line, err := s.logsAdapter.Pop(e.InstanceUUID)
		if err != nil && !errors.Is(err, types.ErrBufferEmpty) {
			log.Error(err)
			return
		}
		if line.Kind == types.LogKindDownloads {
			downloads = line.Message.(*types.LogLineMessageDownloads)
			downloads.Merge(download.DownloadProgress)
		} else {
			downloads = types.NewLogLineMessageDownloads(download.DownloadProgress)
		}
		s.logsAdapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    types.LogKindDownloads,
			Message: downloads,
		})
	default:
		s.logsAdapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
}
