package service

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/types/app"
)

type InstanceLogsService struct {
	uuid    uuid.UUID
	adapter types.InstanceLogsAdapterPort
}

func NewInstanceLogsService(ctx *app.Context, adapter types.InstanceLogsAdapterPort) *InstanceLogsService {
	s := &InstanceLogsService{
		uuid:    uuid.New(),
		adapter: adapter,
	}
	ctx.AddListener(s)
	return s
}

func (s *InstanceLogsService) GetLatestLogs(uuid uuid.UUID) ([]types.LogLine, error) {
	return s.adapter.LoadBuffer(uuid)
}
