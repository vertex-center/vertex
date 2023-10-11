package app

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type AppsRegistry struct {
	uuid uuid.UUID
	ctx  *types.VertexContext

	apps                map[string]App
	mutexApps           *sync.RWMutex
	routers             map[string]Router
	mutexRouters        *sync.RWMutex
	eventListeners      map[string]*types.Listener
	mutexEventListeners *sync.RWMutex
}

func NewAppsRegistry(ctx *types.VertexContext) *AppsRegistry {
	r := &AppsRegistry{
		uuid: uuid.New(),
		ctx:  ctx,

		apps:                map[string]App{},
		mutexApps:           &sync.RWMutex{},
		routers:             map[string]Router{},
		mutexRouters:        &sync.RWMutex{},
		eventListeners:      map[string]*types.Listener{},
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

func (registry *AppsRegistry) RegisterRouter(route string, router Router) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	registry.routers[route] = router
}

func (registry *AppsRegistry) UnregisterRouter(route string) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	delete(registry.routers, route)
}

func (registry *AppsRegistry) GetContext() *types.VertexContext {
	return registry.ctx
}

func (registry *AppsRegistry) GetRouters() map[string]Router {
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
