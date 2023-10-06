package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addMetricsRoutes(r *router.Group) {
	r.GET("/collector/:collector", handleGetMetrics)
	r.POST("/collector/:collector/install", handleInstallMetricsCollector)
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

	err = metricsService.ConfigureInstance(inst)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureMetricsInstance,
			PublicMessage:  "Failed to configure instance to monitor Vertex.",
			PrivateMessage: err.Error(),
		})
		return
	}
}

func handleGetMetrics(c *router.Context) {
	collector, err := getCollector(c)
	if err != nil {
		return
	}

	metrics, err := metricsService.GetMetrics(collector)
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToGetMetrics,
			PublicMessage:  "Failed to get metrics.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(metrics)
}
