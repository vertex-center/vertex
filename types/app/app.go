package app

import (
	"sync"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
)

type App struct {
	id           string
	name         string
	ctx          *Context
	routers      map[string]Router
	routersMutex *sync.RWMutex
}

func New(ctx *types.VertexContext) *App {
	return &App{
		ctx:          NewContext(ctx),
		routers:      map[string]Router{},
		routersMutex: &sync.RWMutex{},
	}
}

func (app *App) Register(appID, appName string) {
	app.id = appID
	app.name = appName
}

func (app *App) RegisterRouter(route string, router Router) {
	app.routersMutex.Lock()
	defer app.routersMutex.Unlock()
	app.routers[route] = router
}

func (app *App) Routers() map[string]Router {
	return app.routers
}

func (app *App) ID() string {
	return app.id
}

func (app *App) Name() string {
	return app.name
}

func (app *App) Context() *Context {
	return app.ctx
}

type Interface interface {
	Initialize(app *App) error
}

type Uninitializable interface {
	Uninitialize() error
}

type Router interface {
	AddRoutes(r *router.Group)
}

type Service interface {
	OnEvent(e interface{})
}

func HeadersSSE(c *router.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
