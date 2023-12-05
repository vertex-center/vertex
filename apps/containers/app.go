package containers

import (
	"github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	"github.com/vertex-center/vertex/apps/containers/handler"
	apptypes "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/router"
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
		ID:                "containers",
		Name:              "Vertex Containers",
		Description:       "Create and manage containers.",
		Icon:              "deployed_code",
		DefaultPort:       "7504",
		DefaultKernelPort: "7505",
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
	container := r.Group("/container/:container_uuid", middleware.Authenticated)
	// docapi:v route /app/containers/container/{container_uuid} vx_containers_get_container
	container.GET("", containerHandler.Get)
	// docapi:v route /app/containers/container/{container_uuid} vx_containers_delete_container
	container.DELETE("", containerHandler.Delete)
	// docapi:v route /app/containers/container/{container_uuid} vx_containers_patch_container
	container.PATCH("", containerHandler.Patch)
	// docapi:v route /app/containers/container/{container_uuid}/start vx_containers_start_container
	container.POST("/start", containerHandler.Start)
	// docapi:v route /app/containers/container/{container_uuid}/stop vx_containers_stop_container
	container.POST("/stop", containerHandler.Stop)
	// docapi:v route /app/containers/container/{container_uuid}/environment vx_containers_patch_environment
	container.PATCH("/environment", containerHandler.PatchEnvironment)
	// docapi:v route /app/containers/container/{container_uuid}/events vx_containers_events_container
	container.GET("/events", apptypes.HeadersSSE, containerHandler.Events)
	// docapi:v route /app/containers/container/{container_uuid}/docker vx_containers_get_docker
	container.GET("/docker", containerHandler.GetDocker)
	// docapi:v route /app/containers/container/{container_uuid}/docker/recreate vx_containers_recreate_docker
	container.POST("/docker/recreate", containerHandler.RecreateDocker)
	// docapi:v route /app/containers/container/{container_uuid}/logs vx_containers_get_logs
	container.GET("/logs", containerHandler.GetLogs)
	// docapi:v route /app/containers/container/{container_uuid}/update/service vx_containers_update_service
	container.POST("/update/service", containerHandler.UpdateService)
	// docapi:v route /app/containers/container/{container_uuid}/versions vx_containers_get_versions
	container.GET("/versions", containerHandler.GetVersions)
	// docapi:v route /app/containers/container/{container_uuid}/wait vx_containers_wait_status
	container.GET("/wait", containerHandler.WaitStatus)

	containersHandler := handler.NewContainersHandler(a.ctx, containerService)
	containers := r.Group("/containers", middleware.Authenticated)
	// docapi:v route /app/containers/containers vx_containers_get_containers
	containers.GET("", containersHandler.Get)
	// docapi:v route /app/containers/containers/tags vx_containers_get_tags
	containers.GET("/tags", containersHandler.GetTags)
	// docapi:v route /app/containers/containers/search vx_containers_search
	containers.GET("/search", containersHandler.Search)
	// docapi:v route /app/containers/containers/checkupdates vx_containers_check_updates
	containers.GET("/checkupdates", containersHandler.CheckForUpdates)
	// docapi:v route /app/containers/containers/events vx_containers_events
	containers.GET("/events", apptypes.HeadersSSE, containersHandler.Events)

	serviceHandler := handler.NewServiceHandler(serviceService, containerService)
	serv := r.Group("/service/:service_id", middleware.Authenticated)
	// docapi:v route /app/containers/service/{service_id} vx_containers_get_service
	serv.GET("", serviceHandler.Get)
	// docapi:v route /app/containers/service/{service_id}/install vx_containers_install_service
	serv.POST("/install", serviceHandler.Install)

	servicesHandler := handler.NewServicesHandler(serviceService)
	services := r.Group("/services")
	// docapi:v route /app/containers/services vx_containers_get_services
	services.GET("", middleware.Authenticated, servicesHandler.Get)
	services.Static("/icons", "./live/services/icons")

	return nil
}

func (a *App) InitializeKernel(r *router.Group) error {
	dockerKernelAdapter = adapter.NewDockerCliAdapter()

	dockerKernelService = service.NewDockerKernelService(dockerKernelAdapter)

	dockerHandler := handler.NewDockerKernelHandler(dockerKernelService)
	docker := r.Group("/docker")
	// docapi:k route /app/containers/docker/containers vx_containers_kernel_get_containers
	docker.GET("/containers", dockerHandler.GetContainers)
	// docapi:k route /app/containers/docker/containers vx_containers_kernel_create_container
	docker.POST("/container", dockerHandler.CreateContainer)
	// docapi:k route /app/containers/docker/containers/{id} vx_containers_kernel_delete_container
	docker.DELETE("/container/:id", dockerHandler.DeleteContainer)
	// docapi:k route /app/containers/docker/containers/{id}/start vx_containers_kernel_start_container
	docker.POST("/container/:id/start", dockerHandler.StartContainer)
	// docapi:k route /app/containers/docker/containers/{id}/stop vx_containers_kernel_stop_container
	docker.POST("/container/:id/stop", dockerHandler.StopContainer)
	// docapi:k route /app/containers/docker/containers/{id}/info vx_containers_kernel_info_container
	docker.GET("/container/:id/info", dockerHandler.InfoContainer)
	// docapi:k route /app/containers/docker/containers/{id}/logs/stdout vx_containers_kernel_logs_stdout_container
	docker.GET("/container/:id/logs/stdout", dockerHandler.LogsStdoutContainer)
	// docapi:k route /app/containers/docker/containers/{id}/logs/stderr vx_containers_kernel_logs_stderr_container
	docker.GET("/container/:id/logs/stderr", dockerHandler.LogsStderrContainer)
	// docapi:k route /app/containers/docker/containers/{id}/wait/{cond} vx_containers_kernel_wait_container
	docker.GET("/container/:id/wait/:cond", dockerHandler.WaitContainer)
	// docapi:k route /app/containers/docker/containers/mounts/{id} vx_containers_kernel_delete_mounts
	docker.DELETE("/container/:id/mounts", dockerHandler.DeleteMounts)

	// docapi:k route /app/containers/docker/image/{id}/info vx_containers_kernel_info_image
	docker.GET("/image/:id/info", dockerHandler.InfoImage)
	// docapi:k route /app/containers/docker/image/pull vx_containers_kernel_pull_image
	docker.POST("/image/pull", dockerHandler.PullImage)
	// docapi:k route /app/containers/docker/image/build vx_containers_kernel_build_image
	docker.POST("/image/build", dockerHandler.BuildImage)

	return nil
}
