package app

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type AppRegistry struct {
	Interface
	*App
}

type AppsRegistry struct {
	uuid uuid.UUID
	ctx  *types.VertexContext

	apps      map[string]AppRegistry
	mutexApps *sync.RWMutex
}

func NewAppsRegistry(ctx *types.VertexContext) *AppsRegistry {
	return &AppsRegistry{
		uuid: uuid.New(),
		ctx:  ctx,

		apps:      map[string]AppRegistry{},
		mutexApps: &sync.RWMutex{},
	}
}

func (registry *AppsRegistry) RegisterApp(app *App, impl Interface) {
	registry.mutexApps.Lock()
	defer registry.mutexApps.Unlock()
	registry.apps[app.id] = AppRegistry{
		Interface: impl,
		App:       app,
	}
}

func (registry *AppsRegistry) Close() {
	for id, app := range registry.apps {
		if a, ok := app.Interface.(Uninitializable); ok {
			log.Info("uninitializing app", vlog.String("id", id))
			err := a.Uninitialize()
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (registry *AppsRegistry) Apps() map[string]AppRegistry {
	return registry.apps
}
