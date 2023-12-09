package app

import (
	"sync"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/core/types"
)

type AppsRegistry struct {
	uuid uuid.UUID
	ctx  *types.VertexContext

	apps      map[string]Interface
	appsMutex *sync.RWMutex
}

func NewAppsRegistry(ctx *types.VertexContext) *AppsRegistry {
	return &AppsRegistry{
		uuid: uuid.New(),
		ctx:  ctx,

		apps:      map[string]Interface{},
		appsMutex: &sync.RWMutex{},
	}
}

func (registry *AppsRegistry) RegisterApp(app Interface) {
	registry.appsMutex.Lock()
	defer registry.appsMutex.Unlock()
	registry.apps[app.Meta().ID] = app
}

func (registry *AppsRegistry) Apps() map[string]Interface {
	return registry.apps
}
