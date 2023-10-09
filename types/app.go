package types

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppRegistry struct {
	uuid uuid.UUID

	apps         map[string]App
	mutexApps    *sync.RWMutex
	routers      map[string]AppRouter
	mutexRouters *sync.RWMutex
}

func NewAppRegistry() *AppRegistry {
	return &AppRegistry{
		uuid: uuid.New(),

		apps:         map[string]App{},
		mutexApps:    &sync.RWMutex{},
		routers:      map[string]AppRouter{},
		mutexRouters: &sync.RWMutex{},
	}
}

func (registry *AppRegistry) RegisterApp(name string, app App) {
	registry.mutexApps.Lock()
	defer registry.mutexApps.Unlock()
	registry.apps[name] = app
}

func (registry *AppRegistry) UnregisterApp(name string) {
	registry.mutexApps.Lock()
	defer registry.mutexApps.Unlock()
	delete(registry.apps, name)
}

func (registry *AppRegistry) RegisterRouter(route string, router AppRouter) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	registry.routers[route] = router
}

func (registry *AppRegistry) UnregisterRouter(route string) {
	registry.mutexRouters.Lock()
	defer registry.mutexRouters.Unlock()
	delete(registry.routers, route)
}

func (registry *AppRegistry) GetRouters() map[string]AppRouter {
	registry.mutexRouters.RLock()
	defer registry.mutexRouters.RUnlock()
	return registry.routers
}

func (registry *AppRegistry) GetUUID() uuid.UUID {
	return registry.uuid
}

func (registry *AppRegistry) OnEvent(e interface{}) {
	for _, app := range registry.apps {
		app.OnEvent(e)
	}
}

type App interface {
	Initialize(registry *AppRegistry) error
	Uninitialize(registry *AppRegistry) error

	Name() string

	OnEvent(e interface{})
}

type AppRouter interface {
	AddRoutes(r *router.Group)
}

type AppService interface {
	OnEvent(e interface{})
}
