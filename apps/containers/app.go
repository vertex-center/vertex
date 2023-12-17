package containers

import (
	authmiddleware "github.com/vertex-center/vertex/apps/auth/middleware"
	"github.com/vertex-center/vertex/apps/containers/adapter"
	"github.com/vertex-center/vertex/apps/containers/core/port"
	"github.com/vertex-center/vertex/apps/containers/core/service"
	"github.com/vertex-center/vertex/apps/containers/database"
	"github.com/vertex-center/vertex/apps/containers/handler"
	"github.com/vertex-center/vertex/apps/containers/meta"
	"github.com/vertex-center/vertex/apps/monitoring/core/types/metric"
	"github.com/vertex-center/vertex/common/app"
	"github.com/vertex-center/vertex/common/app/appmeta"
	"github.com/vertex-center/vertex/common/middleware"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/wI2L/fizz"
)

var (
	containerService port.ContainerService
	tagsService      port.TagsService
	metricsService   port.MetricsService

	dockerKernelService port.DockerService
)

type App struct {
	ctx *app.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) Load(ctx *app.Context) {
	a.ctx = ctx
}

func (a *App) Meta() appmeta.Meta {
	return meta.Meta
}

func (a *App) Initialize() error {
	db, err := storage.NewDB(storage.DBParams{
		ID:         a.Meta().ID,
		SchemaFunc: database.GetSchema,
		Migrations: database.Migrations,
	})
	if err != nil {
		return err
	}

	var (
		caps       = adapter.NewCapDBAdapter(db)
		ports      = adapter.NewPortDBAdapter(db)
		sysctls    = adapter.NewSysctlDBAdapter(db)
		tags       = adapter.NewTagDBAdapter(db)
		volumes    = adapter.NewVolumeDBAdapter(db)
		containers = adapter.NewContainerDBAdapter(db)
		env        = adapter.NewEnvDBAdapter(db)
		logs       = adapter.NewLogsFSAdapter(nil)
		runner     = adapter.NewRunnerDockerAdapter()
		services   = adapter.NewServiceFSAdapter(nil)
	)

	containerService = service.NewContainerService(a.ctx, caps, containers, env, ports, volumes, tags, sysctls, runner, services, logs)
	tagsService = service.NewTagsService(tags)
	metricsService = service.NewMetricsService(a.ctx)

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(authmiddleware.ReadAuth)

	metric.Serve(r, metricsService)

	var (
		servicesHandler   = handler.NewServicesHandler(containerService)
		serviceHandler    = handler.NewServiceHandler(containerService)
		tagsHandler       = handler.NewTagsHandler(tagsService)
		containersHandler = handler.NewContainersHandler(a.ctx, containerService)
		containerHandler  = handler.NewContainerHandler(a.ctx, containerService)

		container  = r.Group("/container/:container_id", "Container", "", authmiddleware.Authenticated)
		containers = r.Group("/containers", "Containers", "", authmiddleware.Authenticated)
		tag        = r.Group("/tag", "Tag", "", authmiddleware.Authenticated)
		tags       = r.Group("/tags", "Tags", "", authmiddleware.Authenticated)
		serv       = r.Group("/service/:service_id", "Service", "", authmiddleware.Authenticated)
		services   = r.Group("/services", "Services", "")
	)

	// Container

	container.GET("", []fizz.OperationOption{
		fizz.ID("getContainer"),
		fizz.Summary("Get a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
	}, containerHandler.Get())

	container.DELETE("", []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "container still running"}),
	}, containerHandler.Delete())

	container.PATCH("", []fizz.OperationOption{
		fizz.ID("patchContainer"),
		fizz.Summary("Patch a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch container"}),
	}, containerHandler.Patch())

	container.POST("/start", []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to start container"}),
	}, containerHandler.Start())

	container.POST("/stop", []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to stop container"}),
	}, containerHandler.Stop())

	container.PUT("/tag/:tag_id", []fizz.OperationOption{
		fizz.ID("addContainerTag"),
		fizz.Summary("Link tag to container"),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to link tag to container"}),
	}, containerHandler.AddContainerTag())

	container.GET("/environment", []fizz.OperationOption{
		fizz.ID("getContainerEnvironment"),
		fizz.Summary("Get container environment"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container environment"}),
	}, containerHandler.GetContainerEnv())

	container.PATCH("/environment", []fizz.OperationOption{
		fizz.ID("patchContainerEnvironment"),
		fizz.Summary("Patch a container environment"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch container environment"}),
	}, containerHandler.PatchEnvironment())

	container.GET("/events", []fizz.OperationOption{
		fizz.ID("eventsContainer"),
		fizz.Summary("Get container events"),
		fizz.Description("Get events for a container, sent as Server-Sent Events (SSE)."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container events"}),
	}, middleware.SSE, containerHandler.Events())

	container.GET("/docker", []fizz.OperationOption{
		fizz.ID("getDockerContainer"),
		fizz.Summary("Get Docker container info"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get Docker container info"}),
	}, containerHandler.GetDocker())

	container.POST("/docker/recreate", []fizz.OperationOption{
		fizz.ID("recreateDockerContainer"),
		fizz.Summary("Recreate Docker container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to recreate Docker container"}),
	}, containerHandler.RecreateDocker())

	container.GET("/logs", []fizz.OperationOption{
		fizz.ID("getContainerLogs"),
		fizz.Summary("Get container logs"),
		fizz.Description("Get latest container logs."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container logs"}),
	}, containerHandler.GetLogs())

	container.GET("/versions", []fizz.OperationOption{
		fizz.ID("getContainerVersions"),
		fizz.Summary("Get container image versions"),
		fizz.Description("Get the possible versions of the container image that can be used."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container image versions"}),
	}, containerHandler.GetVersions())

	container.GET("/wait", []fizz.OperationOption{
		fizz.ID("waitContainerStatus"),
		fizz.Summary("Wait for a status change"),
		fizz.Description("Wait for a status change of the container."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("408", "Timeout", nil, nil, map[string]interface{}{"error": "wait status timeout"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to wait for status change"}),
	}, containerHandler.WaitStatus())

	// Containers

	containers.GET("", []fizz.OperationOption{
		fizz.ID("getContainers"),
		fizz.Summary("Get containers"),
	}, containersHandler.GetContainers())

	containers.GET("/checkupdates", []fizz.OperationOption{
		fizz.ID("checkForUpdates"),
		fizz.Summary("Check for updates"),
	}, containersHandler.CheckForUpdates())

	containers.GET("/events", []fizz.OperationOption{
		fizz.ID("events"),
		fizz.Summary("Get events"),
	}, middleware.SSE, containersHandler.Events())

	// Tags

	tag.GET("", []fizz.OperationOption{
		fizz.ID("getTag"),
		fizz.Summary("Get tag"),
	}, tagsHandler.GetTag())

	tag.POST("", []fizz.OperationOption{
		fizz.ID("createTag"),
		fizz.Summary("Create tag"),
	}, tagsHandler.CreateTag())

	tag.DELETE("/:id", []fizz.OperationOption{
		fizz.ID("deleteTag"),
		fizz.Summary("Delete tag"),
	}, tagsHandler.DeleteTag())

	tags.GET("", []fizz.OperationOption{
		fizz.ID("getTags"),
		fizz.Summary("Get tags"),
	}, tagsHandler.GetTags())

	// Service

	serv.GET("", []fizz.OperationOption{
		fizz.ID("getService"),
		fizz.Summary("Get service"),
	}, serviceHandler.Get())

	serv.POST("/install", []fizz.OperationOption{
		fizz.ID("installService"),
		fizz.Summary("Install service"),
	}, serviceHandler.Install())

	// Services

	services.GET("", []fizz.OperationOption{
		fizz.ID("getServices"),
		fizz.Summary("Get services"),
	}, authmiddleware.Authenticated, servicesHandler.Get())

	services.GinRouterGroup().Static("/icons", "./live/services/icons")

	return nil
}

func (a *App) InitializeKernel() error {
	dockerKernelAdapter := adapter.NewDockerCliAdapter()
	dockerKernelService = service.NewDockerKernelService(dockerKernelAdapter)
	return nil
}

func (a *App) InitializeKernelRouter(r *fizz.RouterGroup) error {
	var (
		dockerHandler = handler.NewDockerKernelHandler(dockerKernelService)
		docker        = r.Group("/docker", "Docker", "Docker wrapper")
	)

	docker.GET("/containers", []fizz.OperationOption{
		fizz.ID("getContainers"),
		fizz.Summary("Get containers"),
	}, dockerHandler.GetContainers())

	docker.POST("/container", []fizz.OperationOption{
		fizz.ID("createContainer"),
		fizz.Summary("Create container"),
	}, dockerHandler.CreateContainer())

	docker.DELETE("/container/:id", []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete container"),
	}, dockerHandler.DeleteContainer())

	docker.POST("/container/:id/start", []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start container"),
	}, dockerHandler.StartContainer())

	docker.POST("/container/:id/stop", []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop container"),
	}, dockerHandler.StopContainer())

	docker.GET("/container/:id/info", []fizz.OperationOption{
		fizz.ID("infoContainer"),
		fizz.Summary("Get container info"),
	}, dockerHandler.InfoContainer())

	docker.GET("/container/:id/logs/stdout", []fizz.OperationOption{
		fizz.ID("logsStdoutContainer"),
		fizz.Summary("Get container stdout logs"),
		fizz.Description("Get container stdout logs as a stream."),
	}, dockerHandler.LogsStdoutContainer())

	docker.GET("/container/:id/logs/stderr", []fizz.OperationOption{
		fizz.ID("logsStderrContainer"),
		fizz.Summary("Get container stderr logs"),
		fizz.Description("Get container stderr logs as a stream."),
	}, dockerHandler.LogsStderrContainer())

	docker.GET("/container/:id/wait/:cond", []fizz.OperationOption{
		fizz.ID("waitContainer"),
		fizz.Summary("Wait container"),
	}, dockerHandler.WaitContainer())

	docker.DELETE("/container/:id/mounts", []fizz.OperationOption{
		fizz.ID("deleteMounts"),
		fizz.Summary("Delete mounts"),
	}, dockerHandler.DeleteMounts())

	docker.GET("/image/:id/info", []fizz.OperationOption{
		fizz.ID("infoImage"),
		fizz.Summary("Get image info"),
	}, dockerHandler.InfoImage())

	docker.POST("/image/pull", []fizz.OperationOption{
		fizz.ID("pullImage"),
		fizz.Summary("Pull image"),
	}, dockerHandler.PullImage())

	docker.POST("/image/build", []fizz.OperationOption{
		fizz.ID("buildImage"),
		fizz.Summary("Build image"),
	}, dockerHandler.BuildImage())

	return nil
}
