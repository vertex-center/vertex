package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/common/app"
)

type logsService struct {
	uuid    uuid.UUID
	adapter port.LogsAdapter
}

func NewLogsService(ctx *app.Context, adapter port.LogsAdapter) port.LogsService {
	s := &logsService{
		uuid:    uuid.New(),
		adapter: adapter,
	}
	ctx.AddListener(s)
	return s
}

func (s *logsService) GetLatestLogs(uuid types.ContainerID) ([]types.LogLine, error) {
	return s.adapter.LoadBuffer(uuid)
}
