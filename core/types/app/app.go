package app

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

type Meta struct {
	// ID is the unique identifier of the app.
	ID string `json:"id"`

	// Name is the name of the app visible to the user.
	Name string `json:"name"`

	// Description is a brief description of the app.
	Description string `json:"description"`

	// Icon is the material symbol name for the app.
	Icon string `json:"icon"`

	// Category is the category of the app.
	Category string `json:"category"`

	// DefaultPort is the default port of the app.
	DefaultPort string `json:"port"`

	// DefaultKernelPort is the default port of the app in kernel mode.
	DefaultKernelPort string `json:"kernel_port"`

	// Hidden is a flag that indicates if the app does only backend work and should be hidden from the frontend.
	Hidden bool `json:"hidden"`

	// Dependencies is a list of app IDs that this app depends on.
	Dependencies []*Meta `json:"-"`
}

type Interface interface {
	Load(ctx *Context)
	Meta() Meta
}

type InterfaceKernel interface {
	LoadKernel(ctx *Context)
	Meta() Meta
}

type Initializable interface {
	Interface
	Initialize(r *router.Group) error
}

type KernelInitializable interface {
	Interface
	InitializeKernel(r *router.Group) error
}

type Uninitializable interface {
	Interface
	Uninitialize() error
}

type KernelUninitializable interface {
	Interface
	UninitializeKernel() error
}

type DependenciesProvider interface {
	DownloadDependencies() error
}

type HttpHandler func(r *router.Group)

type Service interface {
	OnEvent(e event.Event) error
}

func HeadersSSE(c *router.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

func (a Meta) ApiURL() *url.URL {
	return config.Current.URL(a.ID)
}

func (a Meta) ApiKernelURL() *url.URL {
	return config.Current.KernelURL(a.ID)
}

func (a Meta) DefaultApiURL() string {
	return fmt.Sprintf(config.DefaultApiURLFormat, config.Current.LocalIP(), a.DefaultPort)
}

func (a Meta) DefaultApiKernelURL() string {
	return fmt.Sprintf(config.DefaultApiURLFormat, config.Current.LocalIP(), a.DefaultKernelPort)
}

// RunStandalone runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandalone(app Interface) {
	vertexCtx := types.NewVertexContext(types.About{}, false)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)

	u := app.Meta().ApiURL()

	srv := server.New(app.Meta().ID, u, vertexCtx)

	if a, ok := app.(Initializable); ok {
		base := srv.Router.Group(u.Path)

		err := a.Initialize(base)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		base.GET("/ping", func(c *router.Context) {
			c.OK()
		})
	}

	_ = srv.StartAsync()

	if a, ok := app.(Uninitializable); ok {
		err := a.Uninitialize()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

// RunStandaloneKernel runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandaloneKernel(app Interface) {
	vertexCtx := types.NewVertexContext(types.About{}, true)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)

	u := app.Meta().ApiKernelURL()

	srv := server.New(app.Meta().ID, u, vertexCtx)

	if a, ok := app.(KernelInitializable); ok {
		base := srv.Router.Group(u.Path)

		err := a.InitializeKernel(base)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		base.GET("/ping", func(c *router.Context) {
			c.OK()
		})
	}

	_ = srv.StartAsync()

	if a, ok := app.(KernelUninitializable); ok {
		err := a.UninitializeKernel()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}
