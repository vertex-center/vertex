package app

import (
	"sync"

	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
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
}

type App struct {
	meta         Meta
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

func (app *App) Register(meta Meta) {
	app.meta = meta
}

func (app *App) RegisterRouter(route string, router Router) {
	app.routersMutex.Lock()
	defer app.routersMutex.Unlock()
	app.routers[route] = router
}

func (app *App) Routers() map[string]Router {
	return app.routers
}

func (app *App) Meta() Meta {
	return app.meta
}

func (app *App) ID() string {
	return app.meta.ID
}

func (app *App) Name() string {
	return app.meta.Name
}

func (app *App) Description() string {
	return app.meta.Description
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
