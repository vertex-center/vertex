package containers

import (
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	"github.com/vertex-center/vertex/apps/containers/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

const (
	AppRoute = "/vx-containers"
)

var (
	containerAdapter         port.ContainerAdapter
	containerEnvAdapter      port.ContainerEnvAdapter
	containerLogsAdapter     port.ContainerLogsAdapter
	containerRunnerAdapter   port.ContainerRunnerAdapter
	containerServiceAdapter  port.ContainerServiceAdapter
	containerSettingsAdapter port.ContainerSettingsAdapter

	containerService         port.ContainerService
	containerEnvService      port.ContainerEnvService
	containerLogsService     port.ContainerLogsService
	containerRunnerService   port.ContainerRunnerService
	containerServiceService  port.ContainerServiceService
	containerSettingsService port.ContainerSettingsService
	metricsService           port.MetricsService
	serviceService           port.ServiceService
)

type App struct {
	*apptypes.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize(app *apptypes.App) error {
	a.App = app

	containerAdapter = adapter.NewContainerFSAdapter(nil)
	containerEnvAdapter = adapter.NewContainerEnvFSAdapter(nil)
	containerLogsAdapter = adapter.NewContainerLogsFSAdapter(nil)
	containerRunnerAdapter = adapter.NewContainerRunnerFSAdapter()
	containerServiceAdapter = adapter.NewContainerServiceFSAdapter(nil)
	containerSettingsAdapter = adapter.NewContainerSettingsFSAdapter(nil)

	containerEnvService = service.NewContainerEnvService(containerEnvAdapter)
	containerLogsService = service.NewContainerLogsService(app.Context(), containerLogsAdapter)
	containerRunnerService = service.NewContainerRunnerService(app.Context(), containerRunnerAdapter)
	containerServiceService = service.NewContainerServiceService(containerServiceAdapter)
	containerSettingsService = service.NewContainerSettingsService(containerSettingsAdapter)
	containerService = service.NewContainerService(service.ContainerServiceParams{
		Ctx:                      app.Context(),
		ContainerAdapter:         containerAdapter,
		ContainerRunnerService:   containerRunnerService,
		ContainerServiceService:  containerServiceService,
		ContainerEnvService:      containerEnvService,
		ContainerSettingsService: containerSettingsService,
	})
	metricsService = service.NewMetricsService(app.Context())
	serviceService = service.NewServiceService()

	app.Register(apptypes.Meta{
		ID:          "vx-containers",
		Name:        "Vertex Containers",
		Description: "Create and manage containers.",
		Icon:        "deployed_code",
	})

	app.RegisterRoutes(AppRoute, func(r *router.Group) {
		containerHandler := handler.NewContainerHandler(handler.ContainerHandlerParams{
			Ctx:                      app.Context(),
			ContainerService:         containerService,
			ContainerSettingsService: containerSettingsService,
			ContainerRunnerService:   containerRunnerService,
			ContainerEnvService:      containerEnvService,
			ContainerServiceService:  containerServiceService,
			ContainerLogsService:     containerLogsService,
			ServiceService:           serviceService,
		})
		container := r.Group("/container/:container_uuid")
		container.GET("", containerHandler.Get)
		container.DELETE("", containerHandler.Delete)
		container.PATCH("", containerHandler.Patch)
		container.POST("/start", containerHandler.Start)
		container.POST("/stop", containerHandler.Stop)
		container.PATCH("/environment", containerHandler.PatchEnvironment)
		container.GET("/events", apptypes.HeadersSSE, containerHandler.Events)
		container.GET("/docker", containerHandler.GetDocker)
		container.POST("/docker/recreate", containerHandler.RecreateDocker)
		container.GET("/logs", containerHandler.GetLogs)
		container.POST("/update/service", containerHandler.UpdateService)
		container.GET("/versions", containerHandler.GetVersions)
		container.GET("/wait", containerHandler.Wait)

		containersHandler := handler.NewContainersHandler(app.Context(), containerService)
		containers := r.Group("/containers")
		containers.GET("", containersHandler.Get)
		containers.GET("/tags", containersHandler.GetTags)
		containers.GET("/search", containersHandler.Search)
		containers.GET("/checkupdates", containersHandler.CheckForUpdates)
		containers.GET("/events", apptypes.HeadersSSE, containersHandler.Events)

		serviceHandler := handler.NewServiceHandler(serviceService, containerService)
		serv := r.Group("/service/:service_id")
		serv.GET("", serviceHandler.Get)
		serv.POST("/install", serviceHandler.Install)

		servicesHandler := handler.NewServicesHandler(serviceService)
		services := r.Group("/services")
		services.GET("", servicesHandler.Get)
		services.Static("/icons", "./live/services/icons")
	})

	return nil
}
