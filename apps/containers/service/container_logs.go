package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/core/types/app"
)

type ContainerLogsService struct {
	uuid    uuid.UUID
	adapter types.ContainerLogsAdapterPort
}

func NewContainerLogsService(ctx *app.Context, adapter types.ContainerLogsAdapterPort) *ContainerLogsService {
	s := &ContainerLogsService{
		uuid:    uuid.New(),
		adapter: adapter,
	}
	ctx.AddListener(s)
	return s
}

func (s *ContainerLogsService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.adapter.LoadBuffer(uuid)
}
