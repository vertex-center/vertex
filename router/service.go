package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addServiceRoutes(r *router.Group) {
	r.GET("", handleGetService)
	r.POST("/install", handleServiceInstall)
}

// handleGetService handles the retrieval of a service.
func handleGetService(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(service)
}

// handleServiceInstall handles the installation of a service.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - service_not_found: the service was not found.
//   - failed_to_install_service: failed to install the service.
func handleServiceInstall(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := instanceService.Install(service, "docker")
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
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

	c.JSON(inst)
}
