package app

import (
	"sync"
	"time"

	logsmeta "github.com/vertex-center/vertex/server/apps/logs/meta"
	"github.com/vertex-center/vertex/server/common"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/common/server"
	"github.com/vertex-center/vertex/server/config"
	"github.com/vertex-center/vertex/server/pkg/event"
	"github.com/vertex-center/vertex/server/pkg/net"
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

func RunApps(about common.About, apps []Interface) {
	waitNet()
	wg := sync.WaitGroup{}
	wg.Add(len(apps))
	for _, a := range apps {
		if _, ok := a.(Initializable); !ok {
			continue
		}
		go func(a Interface) {
			runApp(a, about, false)
			wg.Done()
		}(a)
	}
	wg.Wait()
}

func RunKernelApps(about common.About, apps []Interface) {
	waitNet()
	wg := sync.WaitGroup{}
	wg.Add(len(apps))
	for _, a := range apps {
		if _, ok := a.(KernelInitializable); !ok {
			continue
		}
		go func(a Interface) {
			runKernelApp(a, about, false)
			wg.Done()
		}(a)
	}
	wg.Wait()
}

// RunStandalone runs the app as a standalone service.
// It parses the configs, loads the app, initializes it and starts the HTTP server.
func RunStandalone(app Interface, about common.About, waitInternet bool) {
	meta := app.Meta()

	for _, m := range meta.Dependencies {
		config.RegisterHost(m.ID, m.DefaultPort)
	}
	if meta.DefaultKernelPort != "" {
		config.RegisterHost(meta.ID+"-kernel", meta.DefaultKernelPort)
	}
	config.ParseArgs(about)
	config.Current.RegisterAPIAddr(meta.ID, config.Current.DefaultApiAddr(config.Current.Port))

	log.SetupAgent(*config.Current.Addr(logsmeta.Meta.ID))

	runApp(app, about, waitInternet)
}

// RunStandaloneKernel runs the app as a standalone service.
// It loads the app, initializes it and starts the HTTP server.
func RunStandaloneKernel(app Interface, about common.About, waitInternet bool) {
	meta := app.Meta()

	for _, m := range meta.Dependencies {
		config.RegisterHost(m.ID, m.DefaultPort)
	}
	config.ParseArgs(about)
	config.Current.RegisterAPIAddr(meta.ID+"-kernel", config.Current.DefaultApiAddr(config.Current.Port))

	log.SetupAgent(*config.Current.Addr(logsmeta.Meta.ID))

	runKernelApp(app, about, waitInternet)
}

func runApp(app Interface, about common.About, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

	vertexCtx := common.NewVertexContext(about, false)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)
	meta := app.Meta()

	u := config.Current.Addr(meta.ID)

	info := openapi.Info{
		Title:       meta.Name,
		Description: meta.Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(meta.ID, &info, u, vertexCtx)

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
	}

	<-srv.StartAsync()

	if a, ok := app.(Uninitializable); ok {
		err := a.Uninitialize()
		if err != nil {
			panic(err)
		}
	}
}

func runKernelApp(app Interface, about common.About, waitInternet bool) {
	if waitInternet {
		waitNet()
	}

	vertexCtx := common.NewVertexContext(about, true)
	ctx := NewContext(vertexCtx)
	app.Load(ctx)
	meta := app.Meta()

	u := config.Current.KernelAddr(meta.ID)

	info := openapi.Info{
		Title:       meta.Name,
		Description: meta.Description,
		Version:     ctx.About().Version,
	}

	srv := server.New(meta.ID, &info, u, vertexCtx)

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
