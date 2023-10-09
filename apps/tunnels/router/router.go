package router

import (
	"errors"
	"fmt"

	instancesapi "github.com/vertex-center/vertex/apps/instances/api"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

type AppRouter struct{}

func NewAppRouter() *AppRouter {
	return &AppRouter{}
}

func (r *AppRouter) GetServices() []types.AppService {
	return []types.AppService{}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.POST("/provider/:provider/install", r.handleInstallTunnelProvider)
}

func (r *AppRouter) handleInstallTunnelProvider(c *router.Context) {
	provider, err := getTunnelProvider(c)
	if err != nil {
		return
	}

	serv, apiError := instancesapi.GetService(c, provider)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := instancesapi.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	apiError = instancesapi.PatchInstance(c, inst.UUID, types.InstanceSettings{
		Tags: []string{"vertex-cloudflare-tunnel"},
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
			Code:           api.ErrCollectorNotFound,
			PublicMessage:  fmt.Sprintf("Provider not found: %s.", provider),
			PrivateMessage: "The provider is not supported. It should be 'cloudflared'.",
		})
		return "", errors.New("collector not found")
	}
	return provider, nil
}
