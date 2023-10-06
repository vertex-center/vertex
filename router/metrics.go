package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addMetricsRoutes(r *router.Group) {
	r.GET("", handleGetMetrics)
	r.POST("/collector/:collector/install", handleInstallMetricsCollector)
	r.POST("/visualizer/:visualizer/install", handleInstallMetricsVisualizer)
}

func getCollector(c *router.Context) (string, error) {
	collector := c.Param("collector")
	if collector != "prometheus" {
		c.NotFound(router.Error{
			Code:           api.ErrCollectorNotFound,
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
			Code:           api.ErrVisualizerNotFound,
			PublicMessage:  fmt.Sprintf("Visualizer not found: %s.", visualizer),
			PrivateMessage: "The visualizer is not supported. It should be 'grafana'.",
		})
		return "", errors.New("visualizer not found")
	}
	return visualizer, nil
}

func handleInstallMetricsCollector(c *router.Context) {
	collector, err := getCollector(c)
	if err != nil {
		return
	}

	service, err := serviceService.GetById(collector)
	if err != nil {
		return
	}

	inst, err := instanceService.Install(service, "docker")
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", collector),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = metricsService.ConfigureCollector(inst)
	if err == nil {
		err = instanceSettingsService.SetTags(inst, []string{"vertex-prometheus-collector"})
	}
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureMetricsInstance,
			PublicMessage:  "Failed to configure instance to monitor Vertex.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func handleInstallMetricsVisualizer(c *router.Context) {
	visualizer, err := getVisualizer(c)
	if err != nil {
		return
	}

	service, err := serviceService.GetById(visualizer)
	if err != nil {
		return
	}

	inst, err := instanceService.Install(service, "docker")
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", visualizer),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = metricsService.ConfigureVisualizer(inst)
	if err == nil {
		err = instanceSettingsService.SetTags(inst, []string{"vertex-grafana-visualizer"})
	}
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureMetricsInstance,
			PublicMessage:  "Failed to configure instance to monitor Vertex.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func handleGetMetrics(c *router.Context) {
	c.JSON(metricsService.GetMetrics())
}
