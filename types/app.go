package types

import (
	"errors"
	"net/http"
	"sync"

	"github.com/carlmjohnson/requests"
	"github.com/gin-contrib/sse"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

type AppsRegistry struct {
	uuid uuid.UUID
	ctx  *VertexContext

	apps                map[string]App
	mutexApps           *sync.RWMutex
	routers             map[string]AppRouter
	mutexRouters        *sync.RWMutex
	eventListeners      map[string]*Listener
	mutexEventListeners *sync.RWMutex
}

func NewAppsRegistry(ctx *VertexContext) *AppsRegistry {
	r := &AppsRegistry{
		uuid: uuid.New(),
		ctx:  ctx,

		apps:                map[string]App{},
		mutexApps:           &sync.RWMutex{},
		routers:             map[string]AppRouter{},
		mutexRouters:        &sync.RWMutex{},
		eventListeners:      map[string]*Listener{},
		mutexEventListeners: &sync.RWMutex{},
	}
	r.ctx.AddListener(r)
	return r
}

func (registry *AppsRegistry) RegisterApp(name string, app App) {
	registry.mutexApps.Lock()
	defer registry.mutexApps.Unlock()
	registry.apps[name] = app
}

func (registry *AppsRegistry) UnregisterApp(name string) {
	registry.mutexApps.Lock()
	defer registry.mutexApps.Unlock()
	delete(registry.apps, name)
}

func (registry *AppsRegistry) RegisterRouter(route string, router AppRouter) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	registry.routers[route] = router
}

func (registry *AppsRegistry) UnregisterRouter(route string) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	delete(registry.routers, route)
}

func (registry *AppsRegistry) GetContext() *VertexContext {
	return registry.ctx
}

func (registry *AppsRegistry) GetRouters() map[string]AppRouter {
	registry.mutexRouters.RLock()
	defer registry.mutexRouters.RUnlock()
	return registry.routers
}

func (registry *AppsRegistry) GetUUID() uuid.UUID {
	return registry.uuid
}

func (registry *AppsRegistry) OnEvent(e interface{}) {
	for _, app := range registry.apps {
		app.OnEvent(e)
	}
}

type App interface {
	Initialize(registry *AppsRegistry) error
	Uninitialize(registry *AppsRegistry) error

	Name() string

	OnEvent(e interface{})
}

type AppRouter interface {
	AddRoutes(r *router.Group)
}

type AppService interface {
	OnEvent(e interface{})
}

type AppApiError struct {
	HttpCode int
	Code     router.ErrCode `json:"code"`
	Message  string         `json:"message"`
}

func (e *AppApiError) RouterError() router.Error {
	return router.Error{
		Code:          e.Code,
		PublicMessage: e.Message,
	}
}

func HandleError(requestError error, apiError AppApiError) *AppApiError {
	if errors.Is(requestError, requests.ErrInvalidHandled) {
		if requests.HasStatusErr(requestError, http.StatusNotFound) {
			apiError.HttpCode = http.StatusNotFound
		} else if requests.HasStatusErr(requestError, http.StatusInternalServerError) {
			apiError.HttpCode = http.StatusInternalServerError
		}
		return &apiError
	} else if requestError != nil {
		return &AppApiError{
			HttpCode: http.StatusInternalServerError,
			Code:     api.ErrInternalError,
			Message:  "Internal error.",
		}
	}
	return nil
}

func HeadersSSE(c *router.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
