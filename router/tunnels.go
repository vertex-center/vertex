package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

func addTunnelsRoutes(r *router.Group) {
	r.POST("/provider/:provider/install", handleInstallTunnelProvider)
}

func getTunnelProvider(c *router.Context) (string, error) {
	provider := c.Param("provider")
	if provider != "cloudflared" {
		c.NotFound(router.Error{
			Code:           api.ErrCollectorNotFound,
			PublicMessage:  fmt.Sprintf("Provider not found: %s.", provider),
			PrivateMessage: "The provider is not supported. It should be 'cloudflared'.",
		})
		return "", errors.New("collector not found")
	}
	return provider, nil
}

func handleInstallTunnelProvider(c *router.Context) {
	provider, err := getTunnelProvider(c)
	if err != nil {
		return
	}

	service, err := serviceService.GetById(provider)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", provider),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := instanceService.Install(service, "docker")
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = instanceSettingsService.SetTags(inst, []string{"vertex-cloudflare-tunnel"})
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureTunnelInstance,
			PublicMessage:  fmt.Sprintf("Failed to configure tunnel '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
