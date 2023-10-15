package service

import (
	"errors"
	"github.com/google/uuid"
	types2 "github.com/vertex-center/vertex/core/types"
	app2 "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type AppsService struct {
	uuid     uuid.UUID
	ctx      *types2.VertexContext
	apps     []app2.Interface
	registry *app2.AppsRegistry
	router   *router.Router
}

func NewAppsService(ctx *types2.VertexContext, router *router.Router, apps []app2.Interface) *AppsService {
	s := &AppsService{
		uuid:     uuid.New(),
		ctx:      ctx,
		apps:     apps,
		registry: app2.NewAppsRegistry(ctx),
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
	case types2.EventServerStart:
		s.StartApps()
	case types2.EventServerStop:
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

func (s *AppsService) startApp(impl app2.Interface) error {
	a := app2.New(s.ctx)
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

func (s *AppsService) All() []app2.Meta {
	var apps []app2.Meta
	for _, a := range s.registry.Apps() {
		apps = append(apps, a.App.Meta())
	}
	return apps
}
