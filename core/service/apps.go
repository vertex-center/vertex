package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type AppsService struct {
	uuid     uuid.UUID
	kernel   bool
	ctx      *types.VertexContext
	apps     []app.Interface
	registry *app.AppsRegistry
	router   *router.Router
}

func NewAppsService(ctx *types.VertexContext, kernel bool, router *router.Router, apps []app.Interface) port.AppsService {
	s := &AppsService{
		uuid:     uuid.New(),
		kernel:   kernel,
		ctx:      ctx,
		apps:     apps,
		registry: app.NewAppsRegistry(ctx),
		router:   router,
	}
	s.ctx.AddListener(s)
	return s
}

func (s *AppsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *AppsService) OnEvent(e event.Event) error {
	switch e.(type) {
	case types.EventServerStart:
		s.StartApps()
	case types.EventServerStop:
		s.StopApps()
	}
	return nil
}

func (s *AppsService) StartApps() {
	log.Info("initializing apps")

	for i := range s.apps {
		ctx := app.NewContext(s.ctx)
		s.apps[i].Load(ctx)
		s.registry.RegisterApp(s.apps[i])
	}

	for _, a := range s.registry.Apps() {
		id := a.Meta().ID

		var err error
		if s.kernel {
			if a, ok := a.(app.KernelInitializable); ok {
				log.Info("initializing kernel app", vlog.String("id", id))
				group := s.router.Group("/api/app/" + id)
				err = a.InitializeKernel(group)
			}
		} else {
			if a, ok := a.(app.Initializable); ok {
				log.Info("initializing app", vlog.String("id", id))
				group := s.router.Group("/api/app/"+id, middleware.ReadAuth)
				err = a.Initialize(group)
			}
		}
		if err != nil {
			log.Error(err)
			log.Error(errors.New("failed to initialize app"),
				vlog.String("id", id))
		}

		s.ctx.DispatchEvent(types.EventAppReady{
			AppID: a.Meta().ID,
		})
	}

	s.ctx.DispatchEvent(types.EventAllAppsReady{})

	log.Info("apps initialized")
}

func (s *AppsService) StopApps() {
	log.Info("uninitializing apps")

	for _, a := range s.registry.Apps() {
		id := a.Meta().ID

		var err error
		if s.kernel {
			if a, ok := a.(app.KernelUninitializable); ok {
				log.Info("uninitializing kernel app", vlog.String("id", id))
				err = a.UninitializeKernel()
			}
		} else {
			if a, ok := a.(app.Uninitializable); ok {
				log.Info("uninitializing app", vlog.String("id", id))
				err = a.Uninitialize()
			}
		}

		if err != nil {
			log.Error(err)
			log.Error(errors.New("failed to uninitialize app"),
				vlog.String("id", id))
		}
	}

	log.Info("apps uninitialized")
}

func (s *AppsService) All() []app.Meta {
	var apps []app.Meta
	for _, a := range s.registry.Apps() {
		apps = append(apps, a.Meta())
	}
	return apps
}
