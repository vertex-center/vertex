package router

import (
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	service2 "github.com/vertex-center/vertex/apps/containers/core/service"
	app2 "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppRouter struct {
	ctx *app2.Context

	containerAdapter         port.ContainerAdapter
	containerEnvAdapter      port.ContainerEnvAdapter
	containerLogsAdapter     port.ContainerLogsAdapter
	containerRunnerAdapter   port.ContainerRunnerAdapter
	containerServiceAdapter  port.ContainerServiceAdapter
	containerSettingsAdapter port.ContainerSettingsAdapter

	containerService         *service2.ContainerService
	containerEnvService      *service2.ContainerEnvService
	containerLogsService     *service2.ContainerLogsService
	containerRunnerService   *service2.ContainerRunnerService
	containerServiceService  *service2.ContainerServiceService
	containerSettingsService *service2.ContainerSettingsService

	serviceService *service2.ServiceService
}

func NewAppRouter(ctx *app2.Context) *AppRouter {
	r := &AppRouter{
		ctx:                      ctx,
		containerAdapter:         adapter.NewContainerFSAdapter(nil),
		containerEnvAdapter:      adapter.NewContainerEnvFSAdapter(nil),
		containerLogsAdapter:     adapter.NewContainerLogsFSAdapter(nil),
		containerRunnerAdapter:   adapter.NewContainerRunnerFSAdapter(),
		containerServiceAdapter:  adapter.NewContainerServiceFSAdapter(nil),
		containerSettingsAdapter: adapter.NewContainerSettingsFSAdapter(nil),
	}

	r.serviceService = service2.NewServiceService()
	r.containerEnvService = service2.NewContainerEnvService(r.containerEnvAdapter)
	r.containerLogsService = service2.NewContainerLogsService(ctx, r.containerLogsAdapter)
	r.containerRunnerService = service2.NewContainerRunnerService(ctx, r.containerRunnerAdapter)
	r.containerServiceService = service2.NewContainerServiceService(r.containerServiceAdapter)
	r.containerSettingsService = service2.NewContainerSettingsService(r.containerSettingsAdapter)
	r.containerService = service2.NewContainerService(service2.ContainerServiceParams{
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
	container.GET("/events", app2.HeadersSSE, r.handleContainerEvents)
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
	containers.GET("/events", app2.HeadersSSE, r.handleContainersEvents)

	serv := group.Group("/service/:service_id")
	serv.GET("", r.handleGetService)
	serv.POST("/install", r.handleServiceInstall)

	services := group.Group("/services")
	services.GET("", r.handleGetServices)
	services.Static("/icons", "./live/services/icons")
}
