package handler

import (
	"github.com/gin-gonic/gin"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
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
		client := containersapi.NewContainersClient(c)

		serv, err := client.GetService(c, params.Provider)
		if err != nil {
			return err
		}

		inst, err := client.InstallService(c, serv.ID)
		if err != nil {
			return err
		}

		return client.PatchContainer(c, inst.ID, map[string]interface{}{
			"tags": []string{"Vertex Tunnels", "Vertex Tunnels - Cloudflare"},
		})
	})
}
