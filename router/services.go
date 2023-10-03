package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

func addServicesRoutes(r *gin.RouterGroup) {
	r.GET("/available", handleServicesAvailable)
	r.POST("/install", handleServiceInstall)
	r.Static("/icons", "./live/services/icons")
}

// handleServicesAvailable handles the retrieval of all available services.
func handleServicesAvailable(c *gin.Context) {
	c.JSON(http.StatusOK, serviceService.GetAll())
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
func handleServiceInstall(c *gin.Context) {
	var body downloadBody
	err := c.BindJSON(&body)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrFailedToParseBody,
			Message: fmt.Sprintf("failed to parse request body: %v", err),
		})
		return
	}

	service, err := serviceService.GetById(body.ServiceID)
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrServiceNotFound,
			Message: fmt.Sprintf("service not found: %v", err),
		})
		return
	}

	instance, err := instanceService.Install(service, body.Method)
	if err != nil && errors.Is(err, types.ErrServiceNotFound) {
		_ = c.AbortWithError(http.StatusBadRequest, api.Error{
			Code:    api.ErrServiceNotFound,
			Message: fmt.Sprintf("service not found: %v", err),
		})
		return
	} else if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, api.Error{
			Code:    api.ErrFailedToInstallService,
			Message: fmt.Sprintf("failed to install service: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"instance": instance,
	})
}
