package service

import "github.com/vertex-center/vertex/apps/logs/core/port"

type logsService struct {
	adapter port.LogsAdapter
}

func NewLogsService(adapter port.LogsAdapter) port.LogsService {
	return &logsService{
		adapter: adapter,
	}
}

func (s *logsService) Push(content string) error {
	return s.adapter.Push(content)
}
