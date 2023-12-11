package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vlog"
)

func (s *containerLogsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *containerLogsService) OnEvent(e event.Event) error {
	switch e := e.(type) {
	case types.EventContainerLog:
		s.onLogReceived(e)
	case types.EventContainerLoaded:
		log.Info("registering container logs", vlog.String("uuid", e.Container.UUID.String()))
		return s.adapter.Register(e.Container.UUID)
	case types.EventContainerDeleted:
		log.Info("unregistering container logs", vlog.String("uuid", e.ContainerUUID.String()))
		return s.adapter.Unregister(e.ContainerUUID)
	case types.EventContainersStopped:
		log.Info("unregistering all container logs")
		return s.adapter.UnregisterAll()
	}
	return nil
}

func (s *containerLogsService) onLogReceived(e types.EventContainerLog) {
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
