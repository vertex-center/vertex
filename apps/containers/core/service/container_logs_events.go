package service

import (
	"errors"
	types2 "github.com/vertex-center/vertex/apps/containers/core/types"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"
)

func (s *ContainerLogsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *ContainerLogsService) OnEvent(e interface{}) {
	switch e := e.(type) {
	case types2.EventContainerLog:
		s.onLogReceived(e)
	case types2.EventContainerLoaded:
		log.Info("registering container logs", vlog.String("uuid", e.Container.UUID.String()))
		err := s.adapter.Register(e.Container.UUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types2.EventContainerDeleted:
		log.Info("unregistering container logs", vlog.String("uuid", e.ContainerUUID.String()))
		err := s.adapter.Unregister(e.ContainerUUID)
		if err != nil {
			log.Error(err)
			return
		}
	case types2.EventContainersStopped:
		log.Info("unregistering all container logs")
		err := s.adapter.UnregisterAll()
		if err != nil {
			log.Error(err)
			return
		}
	}
}

func (s *ContainerLogsService) onLogReceived(e types2.EventContainerLog) {
	switch e.Kind {
	case types2.LogKindDownload:
		var downloads *types2.LogLineMessageDownloads
		download := e.Message.(*types2.LogLineMessageDownload)

		line, err := s.adapter.Pop(e.ContainerUUID)
		if err != nil && !errors.Is(err, types2.ErrBufferEmpty) {
			log.Error(err)
			return
		}
		if line.Kind == types2.LogKindDownloads {
			downloads = line.Message.(*types2.LogLineMessageDownloads)
			downloads.Merge(download.DownloadProgress)
		} else {
			downloads = types2.NewLogLineMessageDownloads(download.DownloadProgress)
		}
		s.adapter.Push(e.ContainerUUID, types2.LogLine{
			Kind:    types2.LogKindDownloads,
			Message: downloads,
		})
	default:
		s.adapter.Push(e.ContainerUUID, types2.LogLine{
			Kind:    e.Kind,
			Message: e.Message,
		})
	}
}
