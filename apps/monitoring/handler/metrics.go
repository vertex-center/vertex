package handler

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type metricsHandler struct {
	metricsService port.MetricsService
}

func NewMetricsHandler(metricsService port.MetricsService) port.MetricsHandler {
	return &metricsHandler{
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
// docapi tags Monitoring
// docapi response 200 {Metrics} The metrics.
// docapi end

func (r *metricsHandler) Get(c *router.Context) {
	c.JSON(r.metricsService.GetMetrics())
}

// docapi begin vx_monitoring_install_collector
// docapi method POST
// docapi summary Install a collector
// docapi tags Monitoring
// docapi query collector {string} The collector to install.
// docapi response 200
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (r *metricsHandler) InstallCollector(c *router.Context) {
	collector, err := getCollector(c)
	if err != nil {
		return
	}

	token := c.MustGet("token").(string)

	err = r.metricsService.InstallCollector(c, token, collector)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToInstallCollector,
			PublicMessage:  "Failed to install collector.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

// docapi begin vx_monitoring_install_visualizer
// docapi method POST
// docapi summary Install a visualizer
// docapi tags Monitoring
// docapi query visualizer {string} The visualizer to install.
// docapi response 200
// docapi response 400
// docapi response 404
// docapi response 500
// docapi end

func (r *metricsHandler) InstallVisualizer(c *router.Context) {
	visualizer, err := getVisualizer(c)
	if err != nil {
		return
	}

	token := c.MustGet("token").(string)

	err = r.metricsService.InstallVisualizer(c, token, visualizer)
	if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToInstallVisualizer,
			PublicMessage:  "Failed to install visualizer.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
