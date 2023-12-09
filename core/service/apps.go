package service

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vlog"

	"github.com/charmbracelet/lipgloss/table"
)

type appsService struct {
	uuid     uuid.UUID
	kernel   bool
	ctx      *types.VertexContext
	apps     []app.Interface
	registry *app.AppsRegistry
}

func NewAppsService(ctx *types.VertexContext, kernel bool, apps []app.Interface) port.AppsService {
	s := &appsService{
		uuid:     uuid.New(),
		kernel:   kernel,
		ctx:      ctx,
		apps:     apps,
		registry: app.NewAppsRegistry(ctx),
	}
	s.ctx.AddListener(s)
	return s
}

func (s *appsService) GetUUID() uuid.UUID {
	return s.uuid
}

func (s *appsService) OnEvent(e event.Event) error {
	switch e.(type) {
	case types.EventServerLoad:
		s.LoadApps()
	case types.EventServerStart:
		s.StartApps()
	case types.EventServerStop:
		s.StopApps()
	}
	return nil
}

func (s *appsService) LoadApps() {
	log.Info("loading apps")

	for i, a := range s.apps {
		ctx := app.NewContext(s.ctx)
		s.apps[i].Load(ctx)
		s.registry.RegisterApp(s.apps[i])

		if _, ok := a.(app.Initializable); ok {
			config.Current.RegisterApiURL(a.Meta().ID, a.Meta().DefaultApiURL())
		}
		if _, ok := a.(app.KernelInitializable); ok {
			config.Current.RegisterKernelApiURL(a.Meta().ID, a.Meta().DefaultApiKernelURL())
		}
	}

	log.Info("apps loaded")
}

func (s *appsService) StartApps() {
	log.Info("initializing apps", vlog.Int("count", len(s.registry.Apps())))

	for _, a := range s.registry.Apps() {
		if s.kernel {
			if a, ok := a.(app.KernelInitializable); ok {
				app.RunStandaloneKernel(a)
			}
		} else {
			if a, ok := a.(app.Initializable); ok {
				app.RunStandalone(a)
			}
		}
		s.ctx.DispatchEvent(types.EventAppReady{
			AppID: a.Meta().ID,
		})
	}

	s.ctx.DispatchEvent(types.EventAllAppsReady{})

	if !s.kernel {
		t := table.New().Headers("App", "API", "Kernel API")
		for _, a := range s.registry.Apps() {
			var (
				apiURL    string
				kernelURL string
			)
			if _, ok := a.(app.Initializable); ok {
				apiURL = a.Meta().ApiURL().String()
			}
			if _, ok := a.(app.KernelInitializable); ok {
				kernelURL = a.Meta().ApiKernelURL().String()
			}
			t.Row(a.Meta().Name, apiURL, kernelURL)
		}
		fmt.Println(t)
	}

	log.Info("apps initialized")
}

func (s *appsService) StopApps() {
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

func (s *appsService) All() []app.Meta {
	var apps []app.Meta
	for _, a := range s.registry.Apps() {
		apps = append(apps, a.Meta())
	}
	return apps
}
