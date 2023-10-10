package router

import (
	"github.com/vertex-center/vertex/apps/instances/adapter"
	"github.com/vertex-center/vertex/apps/instances/service"
	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/router"
	vtypes "github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/app"
)

type AppRouter struct {
	ctx *vtypes.VertexContext

	instanceAdapter         types.InstanceAdapterPort
	instanceEnvAdapter      types.InstanceEnvAdapterPort
	instanceLogsAdapter     types.InstanceLogsAdapterPort
	instanceRunnerAdapter   types.InstanceRunnerAdapterPort
	instanceServiceAdapter  types.InstanceServiceAdapterPort
	instanceSettingsAdapter types.InstanceSettingsAdapterPort

	instanceService         *service.InstanceService
	instanceEnvService      *service.InstanceEnvService
	instanceLogsService     *service.InstanceLogsService
	instanceRunnerService   *service.InstanceRunnerService
	instanceServiceService  *service.InstanceServiceService
	instanceSettingsService *service.InstanceSettingsService

	serviceService *service.ServiceService
}

func NewAppRouter(ctx *vtypes.VertexContext) *AppRouter {
	r := &AppRouter{
		ctx:                     ctx,
		instanceAdapter:         adapter.NewInstanceFSAdapter(nil),
		instanceEnvAdapter:      adapter.NewInstanceEnvFSAdapter(nil),
		instanceLogsAdapter:     adapter.NewInstanceLogsFSAdapter(nil),
		instanceRunnerAdapter:   adapter.NewInstanceRunnerFSAdapter(),
		instanceServiceAdapter:  adapter.NewInstanceServiceFSAdapter(nil),
		instanceSettingsAdapter: adapter.NewInstanceSettingsFSAdapter(nil),
	}

	r.serviceService = service.NewServiceService()
	r.instanceEnvService = service.NewInstanceEnvService(r.instanceEnvAdapter)
	r.instanceLogsService = service.NewInstanceLogsService(r.instanceLogsAdapter)
	r.instanceRunnerService = service.NewInstanceRunnerService(ctx, r.instanceRunnerAdapter)
	r.instanceServiceService = service.NewInstanceServiceService(r.instanceServiceAdapter)
	r.instanceSettingsService = service.NewInstanceSettingsService(r.instanceSettingsAdapter)
	r.instanceService = service.NewInstanceService(service.InstanceServiceParams{
		Ctx:                     ctx,
		InstanceAdapter:         r.instanceAdapter,
		InstanceRunnerService:   r.instanceRunnerService,
		InstanceServiceService:  r.instanceServiceService,
		InstanceEnvService:      r.instanceEnvService,
		InstanceSettingsService: r.instanceSettingsService,
	})

	return r
}

func (r *AppRouter) GetServices() []app.Service {
	return []app.Service{
		r.instanceService,
		r.instanceEnvService,
		r.instanceLogsService,
		r.instanceRunnerService,
		r.instanceSettingsService,
		r.instanceSettingsService,
		r.serviceService,
	}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	instance := group.Group("/instance/:instance_uuid")
	instance.GET("", r.handleGetInstance)
	instance.DELETE("", r.handleDeleteInstance)
	instance.PATCH("", r.handlePatchInstance)
	instance.POST("/start", r.handleStartInstance)
	instance.POST("/stop", r.handleStopInstance)
	instance.PATCH("/environment", r.handlePatchEnvironment)
	instance.GET("/events", app.HeadersSSE, r.handleInstanceEvents)
	instance.GET("/docker", r.handleGetDocker)
	instance.POST("/docker/recreate", r.handleRecreateDockerContainer)
	instance.GET("/logs", r.handleGetLogs)
	instance.POST("/update/service", r.handleUpdateService)
	instance.GET("/versions", r.handleGetVersions)

	instances := group.Group("/instances")
	instances.GET("", r.handleGetInstances)
	instances.GET("/search", r.handleSearchInstances)
	instances.GET("/checkupdates", r.handleCheckForUpdates)
	instances.GET("/events", app.HeadersSSE, r.handleInstancesEvents)

	serv := group.Group("/service/:service_id")
	serv.GET("", r.handleGetService)
	serv.POST("/install", r.handleServiceInstall)

	services := group.Group("/services")
	services.GET("", r.handleGetServices)
	services.Static("/icons", "./live/services/icons")
}
