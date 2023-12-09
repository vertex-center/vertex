package app

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/server"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

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

func RunApps(app []Interface) {
	waitNet()
	for _, a := range app {
		go RunStandalone(a, false)
	}
}

func RunKernelApps(app []Interface) {
	waitNet()
	for _, a := range app {
		go RunStandaloneKernel(a, false)
	}
}

// RunStandalone runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandalone(app Interface, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

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

	<-srv.StartAsync()

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
func RunStandaloneKernel(app Interface, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

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

	<-srv.StartAsync()

	if a, ok := app.(KernelUninitializable); ok {
		err := a.UninitializeKernel()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
}

func waitNet() {
	err := net.WaitInternetConnWithTimeout(20 * time.Second)
	if err != nil {
		log.Error(fmt.Errorf("internet connection not available: %w", err))
		return
	}
}
