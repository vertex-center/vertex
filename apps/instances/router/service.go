package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/pkg/router"
)

// handleGetService handles the retrieval of a service.
func (r *AppRouter) handleGetService(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := r.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
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
func (r *AppRouter) handleServiceInstall(c *router.Context) {
	serviceID := c.Param("service_id")
	if serviceID == "" {
		c.BadRequest(router.Error{
			Code:           types.ErrCodeServiceIdMissing,
			PublicMessage:  "The request was missing the service ID.",
			PrivateMessage: "Field 'service_id' is required.",
		})
		return
	}

	service, err := r.serviceService.GetById(serviceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := r.instanceService.Install(service, "docker")
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           types.ErrCodeServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", serviceID),
			PrivateMessage: err.Error(),
		})
		return
	} else if err != nil {
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(inst)
}
