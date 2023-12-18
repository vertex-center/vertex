package containers

import (
	"os"
	"path"

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
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/common/middleware"
	"github.com/vertex-center/vertex/common/storage"
	"github.com/vertex-center/vertex/common/updater"
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

	if !ctx.Kernel() {
		bl, err := ctx.About().Baseline()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = updater.Install(bl, updater.NewRepositoryUpdater("vertex_services", path.Join(storage.FSPath, "services"), "vertex-center", "services"))
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
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
		servicesHandler   = handler.NewServiceHandler(containerService)
		tagsHandler       = handler.NewTagsHandler(tagsService)
		containersHandler = handler.NewContainerHandler(a.ctx, containerService)

		containers = r.Group("/containers", "Containers", "", authmiddleware.Authenticated)
		tags       = r.Group("/tags", "Tags", "", authmiddleware.Authenticated)
		services   = r.Group("/services", "Services", "")
	)

	// Container

	containers.GET("", []fizz.OperationOption{
		fizz.ID("getContainers"),
		fizz.Summary("Get containers"),
	}, containersHandler.GetContainers())

	containers.GET("/:container_id", []fizz.OperationOption{
		fizz.ID("getContainer"),
		fizz.Summary("Get a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
	}, containersHandler.Get())

	containers.POST("", []fizz.OperationOption{
		fizz.ID("createContainer"),
		fizz.Summary("Create a container"),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to create container"}),
	}, containersHandler.CreateContainer())

	containers.DELETE("/:container_id", []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "container still running"}),
	}, containersHandler.Delete())

	containers.PATCH("/:container_id", []fizz.OperationOption{
		fizz.ID("patchContainer"),
		fizz.Summary("Patch a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch container"}),
	}, containersHandler.Patch())

	containers.POST("/:container_id/start", []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to start container"}),
	}, containersHandler.Start())

	containers.POST("/:container_id/stop", []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop a container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to stop container"}),
	}, containersHandler.Stop())

	containers.PUT("/:container_id/tags/:tag_id", []fizz.OperationOption{
		fizz.ID("addContainerTag"),
		fizz.Summary("Link tag to container"),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to link tag to container"}),
	}, containersHandler.AddContainerTag())

	containers.GET("/:container_id/environment", []fizz.OperationOption{
		fizz.ID("getContainerEnvironment"),
		fizz.Summary("Get container environment"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container environment"}),
	}, containersHandler.GetContainerEnv())

	containers.PATCH("/:container_id/environment", []fizz.OperationOption{
		fizz.ID("patchContainerEnvironment"),
		fizz.Summary("Patch a container environment"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch container environment"}),
	}, containersHandler.PatchEnvironment())

	containers.GET("/:container_id/events", []fizz.OperationOption{
		fizz.ID("eventsContainer"),
		fizz.Summary("Get container events"),
		fizz.Description("Get events for a container, sent as Server-Sent Events (SSE)."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container events"}),
	}, middleware.SSE, containersHandler.ContainerEvents())

	containers.GET("/:container_id/docker", []fizz.OperationOption{
		fizz.ID("getDockerContainer"),
		fizz.Summary("Get Docker container info"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get Docker container info"}),
	}, containersHandler.GetDocker())

	containers.POST("/:container_id/docker/recreate", []fizz.OperationOption{
		fizz.ID("recreateDockerContainer"),
		fizz.Summary("Recreate Docker container"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to recreate Docker container"}),
	}, containersHandler.RecreateDocker())

	containers.GET("/:container_id/logs", []fizz.OperationOption{
		fizz.ID("getContainerLogs"),
		fizz.Summary("Get container logs"),
		fizz.Description("Get latest container logs."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container logs"}),
	}, containersHandler.GetLogs())

	containers.GET("/:container_id/versions", []fizz.OperationOption{
		fizz.ID("getContainerVersions"),
		fizz.Summary("Get container image versions"),
		fizz.Description("Get the possible versions of the container image that can be used."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container image versions"}),
	}, containersHandler.GetVersions())

	containers.GET("/:container_id/wait", []fizz.OperationOption{
		fizz.ID("waitContainerStatus"),
		fizz.Summary("Wait for a status change"),
		fizz.Description("Wait for a status change of the container."),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("408", "Timeout", nil, nil, map[string]interface{}{"error": "wait status timeout"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to wait for status change"}),
	}, containersHandler.WaitStatus())

	containers.GET("/checkupdates", []fizz.OperationOption{
		fizz.ID("checkForUpdates"),
		fizz.Summary("Check for updates"),
	}, containersHandler.CheckForUpdates())

	containers.GET("/events", []fizz.OperationOption{
		fizz.ID("events"),
		fizz.Summary("Get events"),
	}, middleware.SSE, containersHandler.ContainersEvents())

	// Tags

	tags.GET("/:name", []fizz.OperationOption{
		fizz.ID("getTag"),
		fizz.Summary("Get tag"),
	}, tagsHandler.GetTag())

	tags.GET("", []fizz.OperationOption{
		fizz.ID("getTags"),
		fizz.Summary("Get tags"),
	}, tagsHandler.GetTags())

	tags.POST("", []fizz.OperationOption{
		fizz.ID("createTag"),
		fizz.Summary("Create tag"),
	}, tagsHandler.CreateTag())

	tags.DELETE("/:id", []fizz.OperationOption{
		fizz.ID("deleteTag"),
		fizz.Summary("Delete tag"),
	}, tagsHandler.DeleteTag())

	// Services

	services.GET("/:service_id", []fizz.OperationOption{
		fizz.ID("getService"),
		fizz.Summary("Get service"),
	}, servicesHandler.GetService())

	services.GET("", []fizz.OperationOption{
		fizz.ID("getServices"),
		fizz.Summary("Get services"),
	}, authmiddleware.Authenticated, servicesHandler.GetServices())

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

	docker.POST("/containers", []fizz.OperationOption{
		fizz.ID("createContainer"),
		fizz.Summary("Create container"),
	}, dockerHandler.CreateContainer())

	docker.DELETE("/containers/:id", []fizz.OperationOption{
		fizz.ID("deleteContainer"),
		fizz.Summary("Delete container"),
	}, dockerHandler.DeleteContainer())

	docker.POST("/containers/:id/start", []fizz.OperationOption{
		fizz.ID("startContainer"),
		fizz.Summary("Start container"),
	}, dockerHandler.StartContainer())

	docker.POST("/containers/:id/stop", []fizz.OperationOption{
		fizz.ID("stopContainer"),
		fizz.Summary("Stop container"),
	}, dockerHandler.StopContainer())

	docker.GET("/containers/:id/info", []fizz.OperationOption{
		fizz.ID("infoContainer"),
		fizz.Summary("Get container info"),
	}, dockerHandler.InfoContainer())

	docker.GET("/containers/:id/logs/stdout", []fizz.OperationOption{
		fizz.ID("logsStdoutContainer"),
		fizz.Summary("Get container stdout logs"),
		fizz.Description("Get container stdout logs as a stream."),
	}, dockerHandler.LogsStdoutContainer())

	docker.GET("/containers/:id/logs/stderr", []fizz.OperationOption{
		fizz.ID("logsStderrContainer"),
		fizz.Summary("Get container stderr logs"),
		fizz.Description("Get container stderr logs as a stream."),
	}, dockerHandler.LogsStderrContainer())

	docker.GET("/containers/:id/wait/:cond", []fizz.OperationOption{
		fizz.ID("waitContainer"),
		fizz.Summary("Wait container"),
	}, dockerHandler.WaitContainer())

	docker.DELETE("/containers/:id/mounts", []fizz.OperationOption{
		fizz.ID("deleteMounts"),
		fizz.Summary("Delete mounts"),
	}, dockerHandler.DeleteMounts())

	docker.GET("/images/:id/info", []fizz.OperationOption{
		fizz.ID("infoImage"),
		fizz.Summary("Get image info"),
	}, dockerHandler.InfoImage())

	docker.POST("/images/pull", []fizz.OperationOption{
		fizz.ID("pullImage"),
		fizz.Summary("Pull image"),
	}, dockerHandler.PullImage())

	docker.POST("/images/build", []fizz.OperationOption{
		fizz.ID("buildImage"),
		fizz.Summary("Build image"),
	}, dockerHandler.BuildImage())

	return nil
}
