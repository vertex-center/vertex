package containers

import (
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	"github.com/vertex-center/vertex/apps/containers/handler"
	"github.com/vertex-center/vertex/apps/containers/meta"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
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
	return meta.Meta
}

func (a *App) Initialize(r *router.Group) error {
	r.Use(middleware.ReadAuth())

	var (
		containerAdapter         = adapter.NewContainerFSAdapter(nil)
		containerEnvAdapter      = adapter.NewContainerEnvFSAdapter(nil)
		containerLogsAdapter     = adapter.NewContainerLogsFSAdapter(nil)
		containerRunnerAdapter   = adapter.NewContainerRunnerFSAdapter()
		containerServiceAdapter  = adapter.NewContainerServiceFSAdapter(nil)
		containerSettingsAdapter = adapter.NewContainerSettingsFSAdapter(nil)

		serviceService           = service.NewServiceService()
		containerEnvService      = service.NewContainerEnvService(containerEnvAdapter)
		containerLogsService     = service.NewContainerLogsService(a.ctx, containerLogsAdapter)
		containerRunnerService   = service.NewContainerRunnerService(a.ctx, containerRunnerAdapter)
		containerServiceService  = service.NewContainerServiceService(containerServiceAdapter)
		containerSettingsService = service.NewContainerSettingsService(containerSettingsAdapter)
		containerService         = service.NewContainerService(service.ContainerServiceParams{
			Ctx:                      a.ctx,
			ContainerAdapter:         containerAdapter,
			ContainerRunnerService:   containerRunnerService,
			ContainerServiceService:  containerServiceService,
			ContainerEnvService:      containerEnvService,
			ContainerSettingsService: containerSettingsService,
			ServiceService:           serviceService,
		})
		_ = service.NewMetricsService(a.ctx)

		servicesHandler   = handler.NewServicesHandler(serviceService)
		serviceHandler    = handler.NewServiceHandler(serviceService, containerService)
		containersHandler = handler.NewContainersHandler(a.ctx, containerService)
		containerHandler  = handler.NewContainerHandler(handler.ContainerHandlerParams{
			Ctx:                      a.ctx,
			ContainerService:         containerService,
			ContainerSettingsService: containerSettingsService,
			ContainerRunnerService:   containerRunnerService,
			ContainerEnvService:      containerEnvService,
			ContainerServiceService:  containerServiceService,
			ContainerLogsService:     containerLogsService,
			ServiceService:           serviceService,
		})

		container  = r.Group("/container/:container_uuid", "Container", "", middleware.Authenticated())
		containers = r.Group("/containers", "Containers", "", middleware.Authenticated())
		serv       = r.Group("/service/:service_id", "Service", "", middleware.Authenticated())
		services   = r.Group("/services", "Services", "")
	)

	container.GET("", containerHandler.GetInfo(), containerHandler.Get())
	container.DELETE("", containerHandler.DeleteInfo(), containerHandler.Delete())
	container.PATCH("", containerHandler.PatchInfo(), containerHandler.Patch())
	container.POST("/start", containerHandler.StartInfo(), containerHandler.Start())
	container.POST("/stop", containerHandler.StopInfo(), containerHandler.Stop())
	container.PATCH("/environment", containerHandler.PatchEnvironmentInfo(), containerHandler.PatchEnvironment())
	container.GET("/events", containerHandler.EventsInfo(), apptypes.HeadersSSE(), containerHandler.Events())
	container.GET("/docker", containerHandler.GetDockerInfo(), containerHandler.GetDocker())
	container.POST("/docker/recreate", containerHandler.RecreateDockerInfo(), containerHandler.RecreateDocker())
	container.GET("/logs", containerHandler.GetLogsInfo(), containerHandler.GetLogs())
	container.POST("/update/service", containerHandler.UpdateServiceInfo(), containerHandler.UpdateService())
	container.GET("/versions", containerHandler.GetVersionsInfo(), containerHandler.GetVersions())
	container.GET("/wait", containerHandler.WaitStatusInfo(), containerHandler.WaitStatus())

	containers.GET("", containersHandler.GetInfo(), containersHandler.Get())
	containers.GET("/tags", containersHandler.GetTagsInfo(), containersHandler.GetTags())
	containers.GET("/search", containersHandler.SearchInfo(), containersHandler.Search())
	containers.GET("/checkupdates", containersHandler.CheckForUpdatesInfo(), containersHandler.CheckForUpdates())
	containers.GET("/events", containersHandler.EventsInfo(), apptypes.HeadersSSE(), containersHandler.Events())

	serv.GET("", serviceHandler.GetInfo(), serviceHandler.Get())
	serv.POST("/install", serviceHandler.InstallInfo(), serviceHandler.Install())

	services.GET("", servicesHandler.GetInfo(), middleware.Authenticated(), servicesHandler.Get())
	services.GinRouterGroup().Static("/icons", "./live/services/icons")

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	var (
		dockerKernelAdapter = adapter.NewDockerCliAdapter()
		dockerKernelService = service.NewDockerKernelService(dockerKernelAdapter)
		dockerHandler       = handler.NewDockerKernelHandler(dockerKernelService)
		docker              = r.Group("/docker", "Docker", "Docker wrapper")
	)

	docker.GET("/containers", dockerHandler.GetContainersInfo(), dockerHandler.GetContainers())
	docker.POST("/container", dockerHandler.CreateContainerInfo(), dockerHandler.CreateContainer())
	docker.DELETE("/container/:id", dockerHandler.DeleteContainerInfo(), dockerHandler.DeleteContainer())
	docker.POST("/container/:id/start", dockerHandler.StartContainerInfo(), dockerHandler.StartContainer())
	docker.POST("/container/:id/stop", dockerHandler.StopContainerInfo(), dockerHandler.StopContainer())
	docker.GET("/container/:id/info", dockerHandler.InfoContainerInfo(), dockerHandler.InfoContainer())
	docker.GET("/container/:id/logs/stdout", dockerHandler.LogsStdoutContainerInfo(), dockerHandler.LogsStdoutContainer())
	docker.GET("/container/:id/logs/stderr", dockerHandler.LogsStderrContainerInfo(), dockerHandler.LogsStderrContainer())
	docker.GET("/container/:id/wait/:cond", dockerHandler.WaitContainerInfo(), dockerHandler.WaitContainer())
	docker.DELETE("/container/:id/mounts", dockerHandler.DeleteMountsInfo(), dockerHandler.DeleteMounts())

	docker.GET("/image/:id/info", dockerHandler.InfoImageInfo(), dockerHandler.InfoImage())
	docker.POST("/image/pull", dockerHandler.PullImageInfo(), dockerHandler.PullImage())
	docker.POST("/image/build", dockerHandler.BuildImageInfo(), dockerHandler.BuildImage())

	return nil
}
