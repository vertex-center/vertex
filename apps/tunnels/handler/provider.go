package handler

import (
	"github.com/gin-gonic/gin"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/tunnels/core/port"
	"github.com/vertex-center/vertex/pkg/router"
)

type providerHandler struct{}

func NewProviderHandler() port.ProviderHandler {
	return &providerHandler{}
}

type InstallParams struct {
	Provider string `path:"provider"`
}

func (r *providerHandler) Install() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallParams) error {
		token := c.MustGet("token").(string)

		client := containersapi.NewContainersClient(token)

		serv, apiError := client.GetService(c, params.Provider)
		if apiError != nil {
			return apiError.RouterError()
		}

		inst, apiError := client.InstallService(c, serv.ID)
		if apiError != nil {
			return apiError.RouterError()
		}

		apiError = client.PatchContainer(c, inst.UUID, containerstypes.ContainerSettings{
			Tags: []string{"Vertex Tunnels", "Vertex Tunnels - Cloudflare"},
		})
		if apiError != nil {
			return apiError.RouterError()
		}

		return nil
	})
}
