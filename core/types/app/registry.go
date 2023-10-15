package app

import (
	"github.com/vertex-center/vertex/core/types"
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
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
	appsMutex *sync.RWMutex
}

func NewAppsRegistry(ctx *types.VertexContext) *AppsRegistry {
	return &AppsRegistry{
		uuid: uuid.New(),
		ctx:  ctx,

		apps:      map[string]AppRegistry{},
		appsMutex: &sync.RWMutex{},
	}
}

func (registry *AppsRegistry) RegisterApp(app *App, impl Interface) error {
	registry.appsMutex.Lock()
	defer registry.appsMutex.Unlock()
	err := impl.Initialize(app)
	if err != nil {
		return err
	}
	registry.apps[app.ID()] = AppRegistry{
		Interface: impl,
		App:       app,
	}
	return nil
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
