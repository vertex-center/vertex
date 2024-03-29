package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/server/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/server/apps/containers/core/types"
	"github.com/vertex-center/vertex/server/apps/tunnels/core/port"
	"github.com/vertex-center/vertex/server/pkg/router"
)

type providerHandler struct{}

func NewProviderHandler() port.ProviderHandler {
	return &providerHandler{}
}

type InstallParams struct {
	Provider string `path:"provider"`
}

func (r *providerHandler) Install() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InstallParams) error {
		cli := containersapi.NewContainersClient(ctx)

		template, err := cli.GetTemplate(ctx, params.Provider)
		if err != nil {
			return err
		}

		c, err := cli.CreateContainer(ctx, containerstypes.CreateContainerOptions{
			TemplateID: &template.ID,
		})
		if err != nil {
			return err
		}

		tag, err := cli.GetTag(ctx, "Vertex Tunnels")
		if errors.Is(err, errors.NotFound) {
			tag, err = cli.CreateTag(ctx, containerstypes.Tag{
				Name: "Vertex Tunnels",
			})
			if err != nil {
				return err
			}
		} else if err != nil {
			return err
		}

		return cli.AddContainerTag(ctx, c.ID, tag.ID)
	})
}
