package adapter

import (
	"path"
	"sync"

	"github.com/vertex-center/vertex/apps/logs/core/port"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vlog"
)

type FSLogsAdapter struct {
	mu     sync.Mutex
	logger *vlog.Logger
}

func NewFSLogsAdapter() port.LogsAdapter {
	logger := vlog.New(
		vlog.WithOutputFile(vlog.LogFormatJson, path.Join(storage.FSPath, "logs")),
	)
	a := &FSLogsAdapter{
		logger: logger,
	}
	var err error
	if err != nil {
		panic(err)
	}
	return a
}

func (a *FSLogsAdapter) Push(content string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.logger.Raw(content)
	return nil
}
