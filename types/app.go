package types

import (
	"errors"
	"net/http"
	"sync"

	"github.com/carlmjohnson/requests"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
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
