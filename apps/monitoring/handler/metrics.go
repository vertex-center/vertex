package handler

import (
	"errors"
	"fmt"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type MetricsHandler struct {
	metricsService port.MetricsService
}

func NewMetricsHandler(metricsService port.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
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

// docapi begin vx_monitoring_get_metrics
// docapi method GET
// docapi summary Get metrics
// docapi tags Apps/Monitoring
// docapi response 200 {Metrics} The metrics.
// docapi end

func (r *MetricsHandler) Get(c *router.Context) {
	c.JSON(r.metricsService.GetMetrics())
}

// docapi begin vx_monitoring_install_collector
// docapi method POST
// docapi summary Install a collector
// docapi tags Apps/Monitoring
// docapi query collector {string} The collector to install.
// docapi response 200
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (r *MetricsHandler) InstallCollector(c *router.Context) {
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

// docapi begin vx_monitoring_install_visualizer
// docapi method POST
// docapi summary Install a visualizer
// docapi tags Apps/Monitoring
// docapi query visualizer {string} The visualizer to install.
// docapi response 200
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (r *MetricsHandler) InstallVisualizer(c *router.Context) {
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
