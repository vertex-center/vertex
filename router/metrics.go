package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addMetricsRoutes(r *router.Group) {
	r.POST("/install/:collector", handleInstallMetricsCollector)
}

func handleInstallMetricsCollector(c *router.Context) {
	collector := c.Param("collector")
	if collector != "prometheus" {
		c.NotFound(router.Error{
			Code:           api.ErrCollectorNotFound,
			PublicMessage:  fmt.Sprintf("Collector not found: %s.", collector),
			PrivateMessage: "The collector is not supported. It should be 'prometheus'.",
		})
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
