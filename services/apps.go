package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/app"
	"github.com/vertex-center/vlog"
)

type AppsService struct {
	uuid     uuid.UUID
	ctx      *types.VertexContext
	apps     []app.Interface
	registry *app.AppsRegistry
	router   *router.Router
}

func NewAppsService(ctx *types.VertexContext, router *router.Router, apps []app.Interface) *AppsService {
	s := &AppsService{
		uuid:     uuid.New(),
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

func (s *AppsService) OnEvent(e interface{}) {
	switch e.(type) {
	case types.EventServerStart:
		s.StartApps()
	case types.EventServerStop:
		s.StopApps()
	}
}

func (s *AppsService) StartApps() {
	log.Info("starting apps")

	for i := range s.apps {
		err := s.startApp(s.apps[i])
		if err != nil {
			log.Error(errors.New("failed to initialize app"), vlog.String("error", err.Error()))
		}
	}

	for _, a := range s.registry.Apps() {
		for group, rtr := range a.Routers() {
			rtr.AddRoutes(s.router.Group("/api/app" + group))
		}
	}
}

func (s *AppsService) startApp(impl app.Interface) error {
	a := app.New(s.ctx)
	err := s.registry.RegisterApp(a, impl)
	if err != nil {
		log.Error(errors.New("failed to initialize app"), vlog.String("error", err.Error()))
		return err
	}
	log.Info("app initialized", vlog.String("name", a.Name()))
	return nil
}

func (s *AppsService) StopApps() {
	s.registry.Close()
}

func (s *AppsService) All() []app.Meta {
	var apps []app.Meta
	for _, a := range s.registry.Apps() {
		apps = append(apps, a.App.Meta())
	}
	return apps
}
