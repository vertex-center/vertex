package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/core/types/app"
)

type containerLogsService struct {
	uuid    uuid.UUID
	adapter port.ContainerLogsAdapter
}

func NewContainerLogsService(ctx *app.Context, adapter port.ContainerLogsAdapter) port.ContainerLogsService {
	s := &containerLogsService{
		uuid:    uuid.New(),
		adapter: adapter,
	}
	ctx.AddListener(s)
	return s
}

func (s *containerLogsService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.adapter.LoadBuffer(uuid)
}
