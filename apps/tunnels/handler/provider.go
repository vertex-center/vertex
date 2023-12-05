package handler

import (
	"errors"
	"fmt"

	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/monitoring/core/types"
	"github.com/vertex-center/vertex/apps/tunnels/core/port"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/pkg/router"
)

type ProviderHandler struct{}

func NewProviderHandler() port.ProviderHandler {
	return &ProviderHandler{}
}

// docapi begin vx_tunnels_install_provider
// docapi method POST
// docapi summary Install a tunnel provider
// docapi tags Apps/Tunnels
// docapi query provider {string} The provider to install.
// docapi response 204
// docapi response 500
// docapi end

func (r *ProviderHandler) Install(c *router.Context) {
	provider, err := getTunnelProvider(c)
	if err != nil {
		return
	}

	token := c.MustGet("token").(string)

	client := containersapi.NewContainersClient(token)

	serv, apiError := client.GetService(c, provider)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := client.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	apiError = client.PatchContainer(c, inst.UUID, containerstypes.ContainerSettings{
		Tags: []string{"Vertex Tunnels", "Vertex Tunnels - Cloudflare"},
	})
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.OK()
}

func getTunnelProvider(c *router.Context) (string, error) {
	provider := c.Param("provider")
	if provider != "cloudflared" {
		c.NotFound(router.Error{
			Code:           types.ErrCodeCollectorNotFound,
			PublicMessage:  fmt.Sprintf("Provider not found: %s.", provider),
			PrivateMessage: "The provider is not supported. It should be 'cloudflared'.",
		})
		return "", errors.New("collector not found")
	}
	return provider, nil
}
