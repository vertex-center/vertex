package app

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/server"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Interface interface {
	Load(ctx *Context)
	Meta() appmeta.Meta
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

func RunApps(apps []Interface) {
	waitNet()
	for _, a := range apps {
		if _, ok := a.(Initializable); !ok {
			continue
		}

		go RunStandalone(a, false)
	}
}

func RunKernelApps(apps []Interface) {
	waitNet()
	for _, a := range apps {
		if _, ok := a.(KernelInitializable); !ok {
			continue
		}
		go RunStandaloneKernel(a, false)
	}
}

// RunStandalone runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandalone(app Interface, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

	vertexCtx := common.NewVertexContext(common.About{}, false)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)

	config.Current.RegisterApiURL(app.Meta().ID, config.Current.DefaultApiURL(app.Meta().DefaultPort))
	if _, ok := app.(KernelInitializable); ok {
		config.Current.RegisterKernelApiURL(app.Meta().ID, config.Current.DefaultApiURL(app.Meta().DefaultKernelPort))
	}

	u := config.Current.URL(app.Meta().ID)

	info := openapi.Info{
		Title:       app.Meta().Name,
		Description: app.Meta().Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(app.Meta().ID, &info, u, vertexCtx)

	if a, ok := app.(Initializable); ok {
		err := a.Initialize()
		if err != nil {
			panic(err)
		}
	}

	if a, ok := app.(InitializableRouter); ok {
		base := srv.Router.Group(u.Path, "", "")

		err := a.InitializeRouter(base)
		if err != nil {
			panic(err)
		}

		base.GET("/ping", []fizz.OperationOption{
			fizz.ID("ping"),
			fizz.Summary("Ping the app"),
		}, router.Handler(func(c *gin.Context) error {
			return nil
		}))
	}

	<-srv.StartAsync()

	if a, ok := app.(Uninitializable); ok {
		err := a.Uninitialize()
		if err != nil {
			panic(err)
		}
	}
}

// RunStandaloneKernel runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandaloneKernel(app Interface, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

	vertexCtx := common.NewVertexContext(common.About{}, true)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)

	config.Current.RegisterKernelApiURL(app.Meta().ID, config.Current.DefaultApiURL(app.Meta().DefaultKernelPort))

	u := config.Current.KernelURL(app.Meta().ID)

	info := openapi.Info{
		Title:       app.Meta().Name,
		Description: app.Meta().Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(app.Meta().ID, &info, u, vertexCtx)

	if a, ok := app.(KernelInitializable); ok {
		err := a.InitializeKernel()
		if err != nil {
			panic(err)
		}
	}

	if a, ok := app.(KernelInitializableRouter); ok {
		base := srv.Router.Group(u.Path, "", "")

		err := a.InitializeKernelRouter(base)
		if err != nil {
			panic(err)
		}

		base.GET("/ping", []fizz.OperationOption{
			fizz.ID("ping"),
			fizz.Summary("Ping the app"),
		}, router.Handler(func(c *gin.Context) error {
			return nil
		}))
	}

	<-srv.StartAsync()

	if a, ok := app.(KernelUninitializable); ok {
		err := a.UninitializeKernel()
		if err != nil {
			panic(err)
		}
	}
}

func waitNet() {
	err := net.WaitInternetConnWithTimeout(20 * time.Second)
	if err != nil {
		panic(err)
	}
}
