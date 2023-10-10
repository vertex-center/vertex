package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	vtypes "github.com/vertex-center/vertex/types"
)

type InstanceLogsService struct {
	adapter types.InstanceLogsAdapterPort
}

func NewInstanceLogsService(adapter types.InstanceLogsAdapterPort) *InstanceLogsService {
	return &InstanceLogsService{
		adapter: adapter,
	}
}

func (s *InstanceLogsService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.adapter.LoadBuffer(uuid)
}

func (s *InstanceLogsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case vtypes.EventInstanceLog:
		s.onLogReceived(e)
	case vtypes.EventInstanceLoaded:
		err := s.adapter.Register(e.Instance.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	case vtypes.EventInstanceDeleted:
		err := s.adapter.Unregister(e.InstanceUUID)
		if err != nil {
			log.Error(err)
			return
		}
	case vtypes.EventServerStop:
		err := s.adapter.UnregisterAll()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *InstanceLogsService) onLogReceived(e vtypes.EventInstanceLog) {
	switch e.Kind {
	case types.LogKindDownload:
		var downloads *types.LogLineMessageDownloads
		download := e.Message.(*types.LogLineMessageDownload)

		line, err := s.adapter.Pop(e.InstanceUUID)
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
		s.adapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    types.LogKindDownloads,
			Message: downloads,
		})
	default:
		s.adapter.Push(e.InstanceUUID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
}
