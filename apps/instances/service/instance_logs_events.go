package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *InstanceLogsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *InstanceLogsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventInstanceLog:
		s.onLogReceived(e)
	case types.EventInstanceLoaded:
		log.Info("registering instance logs", vlog.String("uuid", e.Instance.UUID.String()))
		err := s.adapter.Register(e.Instance.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types.EventInstanceDeleted:
		log.Info("unregistering instance logs", vlog.String("uuid", e.InstanceUUID.String()))
		err := s.adapter.Unregister(e.InstanceUUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types.EventInstancesStopped:
		log.Info("unregistering all instance logs")
		err := s.adapter.UnregisterAll()
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
