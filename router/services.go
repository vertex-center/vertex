package router

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addServicesRoutes(r *router.Group) {
	r.GET("/available", handleServicesAvailable)
	r.POST("/install", handleServiceInstall)
	r.Static("/icons", "./live/services/icons")
}

// handleServicesAvailable handles the retrieval of all available services.
func handleServicesAvailable(c *router.Context) {
	c.JSON(serviceService.GetAll())
}

type downloadBody struct {
	Method    string `json:"method"`
	ServiceID string `json:"service_id"`
}

// handleServiceInstall handles the installation of a service.
// Errors can be:
//   - failed_to_parse_body: failed to parse the request body.
//   - service_not_found: the service was not found.
//   - failed_to_install_service: failed to install the service.
func handleServiceInstall(c *router.Context) {
	var body downloadBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	service, err := serviceService.GetById(body.ServiceID)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", body.ServiceID),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := instanceService.Install(service, body.Method)
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", body.ServiceID),
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

	c.JSON(gin.H{
		"instance": inst,
	})
}
