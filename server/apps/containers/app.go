package containers

import (
	"os"
	"path"

	authmiddleware "github.com/vertex-center/vertex/server/apps/auth/middleware"
	"github.com/vertex-center/vertex/server/apps/containers/adapter"
	"github.com/vertex-center/vertex/server/apps/containers/core/port"
	"github.com/vertex-center/vertex/server/apps/containers/core/service"
	"github.com/vertex-center/vertex/server/apps/containers/database"
	"github.com/vertex-center/vertex/server/apps/containers/handler"
	"github.com/vertex-center/vertex/server/apps/containers/meta"
	"github.com/vertex-center/vertex/server/apps/monitoring/core/types/metric"
	"github.com/vertex-center/vertex/server/common/app"
	"github.com/vertex-center/vertex/server/common/app/appmeta"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/common/middleware"
	"github.com/vertex-center/vertex/server/common/storage"
	"github.com/vertex-center/vertex/server/common/updater"
	"github.com/wI2L/fizz"
)

var (
	containerService port.ContainerService
	envService       port.EnvService
	tagsService      port.TagsService
	metricsService   port.MetricsService
	portsService     port.PortsService

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
		Name:       a.Meta().ID,
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
		services   = adapter.NewTemplateFSAdapter(nil)
	)

	containerService = service.NewContainerService(a.ctx, caps, containers, env, ports, volumes, tags, sysctls, runner, services, logs)
	envService = service.NewEnvService(env)
	tagsService = service.NewTagsService(tags)
	metricsService = service.NewMetricsService(a.ctx)
	portsService = service.NewPortsService(ports)

	return nil
}

func (a *App) InitializeRouter(r *fizz.RouterGroup) error {
	r.Use(authmiddleware.ReadAuth)

	metric.Serve(r, metricsService)

	var (
		templatesHandler  = handler.NewTemplateHandler(containerService)
		envHandler        = handler.NewEnvHandler(envService)
		tagsHandler       = handler.NewTagsHandler(tagsService)
		containersHandler = handler.NewContainerHandler(a.ctx, containerService)
		portsHandler      = handler.NewPortsHandler(portsService)

		containers   = r.Group("/containers", "Containers", "", authmiddleware.Authenticated)
		environments = r.Group("/environments", "Environment variables", "", authmiddleware.Authenticated)
		ports        = r.Group("/ports", "Ports", "", authmiddleware.Authenticated)
		tags         = r.Group("/tags", "Tags", "", authmiddleware.Authenticated)
		templates    = r.Group("/templates", "Templates", "")
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

	containers.POST("/:container_id/reload", []fizz.OperationOption{
		fizz.ID("reloadContainer"),
		fizz.Summary("Reload a container"),
	}, containersHandler.ReloadContainer())

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

	// Environment

	environments.GET("", []fizz.OperationOption{
		fizz.ID("getEnvironment"),
		fizz.Summary("Get an environment variables"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to get container environment"}),
	}, envHandler.GetEnv())

	environments.PATCH("/:env_id", []fizz.OperationOption{
		fizz.ID("patchEnvironment"),
		fizz.Summary("Patch an environment variable"),
		fizz.Response("404", "Environment variable not found", nil, nil, map[string]interface{}{"error": "environment variable not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch container environment"}),
	}, envHandler.PatchEnv())

	environments.DELETE("/:env_id", []fizz.OperationOption{
		fizz.ID("deleteEnvironment"),
		fizz.Summary("Delete an environment variable"),
		fizz.Response("404", "Environment variable not found", nil, nil, map[string]interface{}{"error": "environment variable not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to delete container environment"}),
	}, envHandler.DeleteEnv())

	environments.POST("", []fizz.OperationOption{
		fizz.ID("createEnvironment"),
		fizz.Summary("Create an environment variable"),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to create container environment"}),
	}, envHandler.CreateEnv())

	// Ports

	ports.GET("", []fizz.OperationOption{
		fizz.ID("getPorts"),
		fizz.Summary("Get ports"),
		fizz.Response("404", "Container not found", nil, nil, map[string]interface{}{"error": "container not found"}),
	}, portsHandler.GetPorts())

	ports.PATCH("/:port_id", []fizz.OperationOption{
		fizz.ID("patchPort"),
		fizz.Summary("Patch ports"),
		fizz.Response("404", "Port not found", nil, nil, map[string]interface{}{"error": "port not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to patch port"}),
	}, portsHandler.PatchPort())

	ports.DELETE("/:port_id", []fizz.OperationOption{
		fizz.ID("deletePort"),
		fizz.Summary("Delete port"),
		fizz.Response("404", "Port not found", nil, nil, map[string]interface{}{"error": "container not found"}),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to delete port"}),
	}, portsHandler.DeletePort())

	ports.POST("", []fizz.OperationOption{
		fizz.ID("createPort"),
		fizz.Summary("Create port"),
		fizz.Response("500", "", nil, nil, map[string]interface{}{"error": "failed to create port"}),
	}, portsHandler.CreatePort())

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

	templates.GET("/:template_id", []fizz.OperationOption{
		fizz.ID("getTemplate"),
		fizz.Summary("Get template"),
	}, templatesHandler.GetTemplate())

	templates.GET("", []fizz.OperationOption{
		fizz.ID("getTemplates"),
		fizz.Summary("Get templates"),
	}, authmiddleware.Authenticated, templatesHandler.GetTemplates())

	templates.GinRouterGroup().Static("/icons", "./live/services/icons")

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

	docker.POST("/volumes", []fizz.OperationOption{
		fizz.ID("createVolume"),
		fizz.Summary("Create volume"),
	}, dockerHandler.CreateVolume())

	docker.DELETE("/volumes", []fizz.OperationOption{
		fizz.ID("deleteVolume"),
		fizz.Summary("Delete volume"),
	}, dockerHandler.DeleteVolume())

	return nil
}
