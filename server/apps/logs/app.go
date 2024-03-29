package logs

import (
	"github.com/vertex-center/vertex/server/apps/logs/adapter"
	"github.com/vertex-center/vertex/server/apps/logs/core/port"
	"github.com/vertex-center/vertex/server/apps/logs/core/service"
	"github.com/vertex-center/vertex/server/apps/logs/handler"
	"github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
	"github.com/wI2L/fizz"
)

var (
	logsService port.LogsService
)

type App struct {
	ctx *app.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx
}

func (a *App) Meta() appmeta.Meta {
	return meta.Meta
}

func (a *App) Initialize() error {
	logsAdapter := adapter.NewFSLogsAdapter()
	logsService = service.NewLogsService(logsAdapter)
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	var (
		logsHandler = handler.NewLogsHandler(logsService)
		logs        = r.Group("/logs", "Logs", "")
	)

	logs.GET("/ws", []fizz.OperationOption{
		fizz.ID("pushLogs"),
		fizz.Summary("Push logs to the server"),
	}, logsHandler.Push())

	return nil
}
