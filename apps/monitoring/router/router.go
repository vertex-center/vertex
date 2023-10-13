package router

import (
	"errors"
	"fmt"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/types"
	"github.com/vertex-center/vertex/apps/monitoring/service"
	"github.com/vertex-center/vertex/apps/monitoring/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/app"
)

type AppRouter struct {
	metricsService *service.MetricsService
}

func NewAppRouter(ctx *app.Context) *AppRouter {
	return &AppRouter{
		metricsService: service.NewMetricsService(ctx),
	}
}

func (r *AppRouter) GetServices() []app.Service {
	return []app.Service{
		r.metricsService,
	}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.GET("/metrics", r.handleGetMetrics)
	group.POST("/collector/:collector/install", r.handleInstallMetricsCollector)
	group.POST("/visualizer/:visualizer/install", r.handleInstallMetricsVisualizer)
}

func getCollector(c *router.Context) (string, error) {
	collector := c.Param("collector")
	if collector != "prometheus" {
		c.NotFound(router.Error{
			Code:           types.ErrCodeCollectorNotFound,
			PublicMessage:  fmt.Sprintf("Collector not found: %s.", collector),
			PrivateMessage: "The collector is not supported. It should be 'prometheus'.",
		})
		return "", errors.New("collector not found")
	}
	return collector, nil
}

func getVisualizer(c *router.Context) (string, error) {
	visualizer := c.Param("visualizer")
	if visualizer != "grafana" {
		c.NotFound(router.Error{
			Code:           types.ErrCodeVisualizerNotFound,
			PublicMessage:  fmt.Sprintf("Visualizer not found: %s.", visualizer),
			PrivateMessage: "The visualizer is not supported. It should be 'grafana'.",
		})
		return "", errors.New("visualizer not found")
	}
	return visualizer, nil
}

func (r *AppRouter) handleGetMetrics(c *router.Context) {
	c.JSON(r.metricsService.GetMetrics())
}

func (r *AppRouter) handleInstallMetricsCollector(c *router.Context) {
	collector, err := getCollector(c)
	if err != nil {
		return
	}

	serv, apiError := containersapi.GetService(c, collector)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := containersapi.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	err = r.metricsService.ConfigureCollector(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToConfigureMetricsContainer,
			PublicMessage:  "Failed to configure container to monitor Vertex.",
			PrivateMessage: err.Error(),
		})
		return
	}

	apiError = containersapi.PatchContainer(c, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Prometheus Collector"},
	})
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.OK()
}

func (r *AppRouter) handleInstallMetricsVisualizer(c *router.Context) {
	visualizer, err := getVisualizer(c)
	if err != nil {
		return
	}

	serv, apiError := containersapi.GetService(c, visualizer)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := containersapi.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	err = r.metricsService.ConfigureVisualizer(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToConfigureMetricsContainer,
			PublicMessage:  "Failed to configure container to monitor Vertex.",
			PrivateMessage: err.Error(),
		})
		return
	}

	apiError = containersapi.PatchContainer(c, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Monitoring", "Vertex Monitoring - Grafana Visualizer"},
	})
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.OK()
}
