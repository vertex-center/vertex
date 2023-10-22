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
	dockerKernelAdapter      port.DockerAdapter

	containerService         port.ContainerService
	containerEnvService      port.ContainerEnvService
	containerLogsService     port.ContainerLogsService
	containerRunnerService   port.ContainerRunnerService
	containerServiceService  port.ContainerServiceService
	containerSettingsService port.ContainerSettingsService
	serviceService           port.ServiceService
	dockerKernelService      port.DockerService
)

type App struct {
	ctx *apptypes.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *apptypes.Context) {
	a.ctx = ctx
}

func (a *App) Meta() apptypes.Meta {
	return apptypes.Meta{
		ID:          "vx-containers",
		Name:        "Vertex Containers",
		Description: "Create and manage containers.",
		Icon:        "deployed_code",
	}
}

func (a *App) Initialize(r *router.Group) error {
	containerAdapter = adapter.NewContainerFSAdapter(nil)
	containerEnvAdapter = adapter.NewContainerEnvFSAdapter(nil)
	containerLogsAdapter = adapter.NewContainerLogsFSAdapter(nil)
	containerRunnerAdapter = adapter.NewContainerRunnerFSAdapter()
	containerServiceAdapter = adapter.NewContainerServiceFSAdapter(nil)
	containerSettingsAdapter = adapter.NewContainerSettingsFSAdapter(nil)

	serviceService = service.NewServiceService()
	containerEnvService = service.NewContainerEnvService(containerEnvAdapter)
	containerLogsService = service.NewContainerLogsService(a.ctx, containerLogsAdapter)
	containerRunnerService = service.NewContainerRunnerService(a.ctx, containerRunnerAdapter)
	containerServiceService = service.NewContainerServiceService(containerServiceAdapter)
	containerSettingsService = service.NewContainerSettingsService(containerSettingsAdapter)
	containerService = service.NewContainerService(service.ContainerServiceParams{
		Ctx:                      a.ctx,
		ContainerAdapter:         containerAdapter,
		ContainerRunnerService:   containerRunnerService,
		ContainerServiceService:  containerServiceService,
		ContainerEnvService:      containerEnvService,
		ContainerSettingsService: containerSettingsService,
		ServiceService:           serviceService,
	})
	service.NewMetricsService(a.ctx)

	containerHandler := handler.NewContainerHandler(handler.ContainerHandlerParams{
		Ctx:                      a.ctx,
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

	containersHandler := handler.NewContainersHandler(a.ctx, containerService)
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

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	dockerKernelAdapter = adapter.NewDockerCliAdapter()

	dockerKernelService = service.NewDockerKernelService(dockerKernelAdapter)

	dockerHandler := handler.NewDockerKernelHandler(dockerKernelService)
	docker := r.Group("/docker")
	docker.GET("/containers", dockerHandler.GetContainers)
	docker.POST("/container", dockerHandler.CreateContainer)
	docker.DELETE("/container/:id", dockerHandler.DeleteContainer)
	docker.POST("/container/:id/start", dockerHandler.StartContainer)
	docker.POST("/container/:id/stop", dockerHandler.StopContainer)
	docker.GET("/container/:id/info", dockerHandler.InfoContainer)
	docker.GET("/container/:id/logs/stdout", dockerHandler.LogsStdoutContainer)
	docker.GET("/container/:id/logs/stderr", dockerHandler.LogsStderrContainer)
	docker.GET("/container/:id/wait/:cond", dockerHandler.WaitContainer)

	docker.GET("/image/:id/info", dockerHandler.InfoImage)
	docker.POST("/image/pull", dockerHandler.PullImage)
	docker.POST("/image/build", dockerHandler.BuildImage)

	docker.POST("/volume", dockerHandler.CreateVolume)
	docker.DELETE("/volume/:name", dockerHandler.DeleteVolume)

	return nil
}
