package serviceeditor

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	logsmeta "github.com/vertex-center/vertex/apps/logs/meta"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/port"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/wI2L/fizz"
)

var (
	editorService port.EditorService
)

var Meta = appmeta.Meta{
	ID:          "devtools-service-editor",
	Name:        "Vertex Service Editor",
	Description: "Create services for publishing.",
	Icon:        "frame_source",
	Category:    "devtools",
	DefaultPort: "7510",
	Dependencies: []*appmeta.Meta{
		&authmeta.Meta,
		&logsmeta.Meta,
	},
}

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
	return Meta
}

func (a *App) Initialize() error {
	editorService = service.NewEditorService()
	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(middleware.ReadAuth)

	editorHandler := handler.NewEditorHandler(editorService)
	editor := r.Group("/editor", "Editor", "Service editor routes", middleware.Authenticated)

	editor.POST("/to-yaml", []fizz.OperationOption{
		fizz.ID("toYaml"),
		fizz.Summary("Convert service to yaml"),
		fizz.Description("Convert service description to a reusable yaml file."),
	}, editorHandler.ToYaml())

	return nil
}
