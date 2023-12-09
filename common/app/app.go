package app

import (
	"fmt"
	"net/url"
	"os"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
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

type Initializable interface {
	Interface
	Initialize() error
}

type InitializableRouter interface {
	Interface
	InitializeRouter(r *fizz.RouterGroup) error
}

type KernelInitializable interface {
	Interface
	InitializeKernel() error
}

type KernelInitializableRouter interface {
	Interface
	InitializeKernelRouter(r *fizz.RouterGroup) error
}

type Uninitializable interface {
	Interface
	Uninitialize() error
}

type KernelUninitializable interface {
	Interface
	UninitializeKernel() error
}

type Service interface {
	OnEvent(e event.Event) error
}

func HeadersSSE() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", sse.ContentType)
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}
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

	info := openapi.Info{
		Title:       app.Meta().Name,
		Description: app.Meta().Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(app.Meta().ID, &info, u, vertexCtx)

	if a, ok := app.(Initializable); ok {
		err := a.Initialize()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	if a, ok := app.(InitializableRouter); ok {
		base := srv.Router.Group(u.Path, "", "")

		err := a.InitializeRouter(base)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		base.GET("/ping", []fizz.OperationOption{
			fizz.Summary("Ping the app"),
		}, router.Handler(func(c *gin.Context) error {
			return nil
		}))
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

	info := openapi.Info{
		Title:       app.Meta().Name,
		Description: app.Meta().Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(app.Meta().ID, &info, u, vertexCtx)

	if a, ok := app.(KernelInitializable); ok {
		err := a.InitializeKernel()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}

	if a, ok := app.(KernelInitializableRouter); ok {
		base := srv.Router.Group(u.Path, "", "")

		err := a.InitializeKernelRouter(base)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		base.GET("/ping", []fizz.OperationOption{
			fizz.Summary("Ping the app"),
		}, router.Handler(func(c *gin.Context) error {
			return nil
		}))
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
