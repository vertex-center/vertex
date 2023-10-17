package serviceeditor

import (
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-devtools-service-editor"
)

type App struct {
	*apptypes.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app

	editorService := service.NewEditorService()

	app.Register(apptypes.Meta{
		ID:          "vx-devtools-service-editor",
		Name:        "Vertex Service Editor",
		Description: "Create services for publishing.",
		Icon:        "frame_source",
	})

	app.RegisterRoutes(AppRoute, func(r *router.Group) {
		editorHandler := handler.NewEditorHandler(editorService)
		editor := r.Group("/editor")
		editor.POST("/to-yaml", editorHandler.ToYaml)
	})

	return nil
}
