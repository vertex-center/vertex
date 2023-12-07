package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vertex-center/vertex/apps/monitoring/core/port"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
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

func (r *metricsHandler) Get(c *router.Context) {
	c.JSON(r.metricsService.GetMetrics())
}

func (r *metricsHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get metrics"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.Metric{}),
		),
	}
}

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

func (r *metricsHandler) InstallCollectorInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Install a collector"),
		oapi.Response(http.StatusOK),
	}
}

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

func (r *metricsHandler) InstallVisualizerInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Install a visualizer"),
		oapi.Response(http.StatusNoContent),
	}
}
