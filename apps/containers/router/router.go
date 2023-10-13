package router

import (
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/service"
	"github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/app"
)

type AppRouter struct {
	ctx *app.Context

	containerAdapter         types.ContainerAdapterPort
	containerEnvAdapter      types.ContainerEnvAdapterPort
	containerLogsAdapter     types.ContainerLogsAdapterPort
	containerRunnerAdapter   types.ContainerRunnerAdapterPort
	containerServiceAdapter  types.ContainerServiceAdapterPort
	containerSettingsAdapter types.ContainerSettingsAdapterPort

	containerService         *service.ContainerService
	containerEnvService      *service.ContainerEnvService
	containerLogsService     *service.ContainerLogsService
	containerRunnerService   *service.ContainerRunnerService
	containerServiceService  *service.ContainerServiceService
	containerSettingsService *service.ContainerSettingsService

	serviceService *service.ServiceService
}

func NewAppRouter(ctx *app.Context) *AppRouter {
	r := &AppRouter{
		ctx:                      ctx,
		containerAdapter:         adapter.NewContainerFSAdapter(nil),
		containerEnvAdapter:      adapter.NewContainerEnvFSAdapter(nil),
		containerLogsAdapter:     adapter.NewContainerLogsFSAdapter(nil),
		containerRunnerAdapter:   adapter.NewContainerRunnerFSAdapter(),
		containerServiceAdapter:  adapter.NewContainerServiceFSAdapter(nil),
		containerSettingsAdapter: adapter.NewContainerSettingsFSAdapter(nil),
	}

	r.serviceService = service.NewServiceService()
	r.containerEnvService = service.NewContainerEnvService(r.containerEnvAdapter)
	r.containerLogsService = service.NewContainerLogsService(ctx, r.containerLogsAdapter)
	r.containerRunnerService = service.NewContainerRunnerService(ctx, r.containerRunnerAdapter)
	r.containerServiceService = service.NewContainerServiceService(r.containerServiceAdapter)
	r.containerSettingsService = service.NewContainerSettingsService(r.containerSettingsAdapter)
	r.containerService = service.NewContainerService(service.ContainerServiceParams{
		Ctx:                      ctx,
		ContainerAdapter:         r.containerAdapter,
		ContainerRunnerService:   r.containerRunnerService,
		ContainerServiceService:  r.containerServiceService,
		ContainerEnvService:      r.containerEnvService,
		ContainerSettingsService: r.containerSettingsService,
	})

	return r
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	container := group.Group("/container/:container_uuid")
	container.GET("", r.handleGetContainer)
	container.DELETE("", r.handleDeleteContainer)
	container.PATCH("", r.handlePatchContainer)
	container.POST("/start", r.handleStartContainer)
	container.POST("/stop", r.handleStopContainer)
	container.PATCH("/environment", r.handlePatchEnvironment)
	container.GET("/events", app.HeadersSSE, r.handleContainerEvents)
	container.GET("/docker", r.handleGetDocker)
	container.POST("/docker/recreate", r.handleRecreateDockerContainer)
	container.GET("/logs", r.handleGetLogs)
	container.POST("/update/service", r.handleUpdateService)
	container.GET("/versions", r.handleGetVersions)
	container.GET("/wait", r.handleWaitContainer)

	containers := group.Group("/containers")
	containers.GET("", r.handleGetContainers)
	containers.GET("/tags", r.handleGetTags)
	containers.GET("/search", r.handleSearchContainers)
	containers.GET("/checkupdates", r.handleCheckForUpdates)
	containers.GET("/events", app.HeadersSSE, r.handleContainersEvents)

	serv := group.Group("/service/:service_id")
	serv.GET("", r.handleGetService)
	serv.POST("/install", r.handleServiceInstall)

	services := group.Group("/services")
	services.GET("", r.handleGetServices)
	services.Static("/icons", "./live/services/icons")
}
