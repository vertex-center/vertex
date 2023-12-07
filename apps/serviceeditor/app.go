package serviceeditor

import (
	authmeta "github.com/vertex-center/vertex/apps/auth/meta"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/serviceeditor/core/service"
	"github.com/vertex-center/vertex/apps/serviceeditor/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

// docapi:service_editor title Vertex Devtools Service Editor
// docapi:service_editor description A service editor for Vertex.
// docapi:service_editor version 0.0.0
// docapi:service_editor filename service_editor

// docapi:service_editor url http://{ip}:{port-kernel}/api
// docapi:service_editor urlvar ip localhost The IP address of the server.
// docapi:service_editor urlvar port-kernel 7510 The port of the server.

var Meta = apptypes.Meta{
	ID:          "devtools-service-editor",
	Name:        "Vertex Service Editor",
	Description: "Create services for publishing.",
	Icon:        "frame_source",
	Category:    "devtools",
	DefaultPort: "7510",
	Dependencies: []*apptypes.Meta{
		&authmeta.Meta,
	},
}

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return Meta
}

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth)

	var (
		editorService = service.NewEditorService()
		editorHandler = handler.NewEditorHandler(editorService)
		editor        = r.Group("/editor", middleware.Authenticated)
	)

	// docapi:service_editor route /editor/to-yaml vx_devtools_service_editor_to_yaml
	editor.POST("/to-yaml", editorHandler.ToYaml)

	return nil
}
