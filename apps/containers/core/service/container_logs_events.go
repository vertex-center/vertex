package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *ContainerLogsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ContainerLogsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types.EventContainerLog:
		s.onLogReceived(e)
	case types.EventContainerLoaded:
		log.Info("registering container logs", vlog.String("uuid", e.Container.UUID.String()))
		err := s.adapter.Register(e.Container.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types.EventContainerDeleted:
		log.Info("unregistering container logs", vlog.String("uuid", e.ContainerUUID.String()))
		err := s.adapter.Unregister(e.ContainerUUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types.EventContainersStopped:
		log.Info("unregistering all container logs")
		err := s.adapter.UnregisterAll()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *ContainerLogsService) onLogReceived(e types.EventContainerLog) {
	switch e.Kind {
	case types.LogKindDownload:
		var downloads *types.LogLineMessageDownloads
		download := e.Message.(*types.LogLineMessageDownload)

		line, err := s.adapter.Pop(e.ContainerUUID)
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
		s.adapter.Push(e.ContainerUUID, types.LogLine{
			Kind:    types.LogKindDownloads,
			Message: downloads,
		})
	default:
		s.adapter.Push(e.ContainerUUID, types.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
}
