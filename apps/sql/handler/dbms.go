package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/pkg/router"
)

type dbmsHandler struct {
	sqlService port.SqlService
}

func NewDBMSHandler(sqlService port.SqlService) port.DBMSHandler {
	return &dbmsHandler{
		sqlService: sqlService,
	}
}

type GetParams struct {
	UUID uuid.NullUUID `path:"container_uuid"`
}

func (r *dbmsHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetParams) (*types.DBMS, error) {
		token := c.MustGet("token").(string)
		client := containersapi.NewContainersClient(token)

		inst, err := client.GetContainer(c, params.UUID.UUID)
		if err != nil {
			return nil, err
		}

		dbms, err := r.sqlService.Get(inst)
		if err != nil {
			return nil, errors.NewNotFound(err, "SQL Database not found")
		}
		return &dbms, nil
	})
}

type InstallParams struct {
	DBMS string `path:"dbms"`
}

func (r *dbmsHandler) Install() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *InstallParams) (*containerstypes.Container, error) {
		token := c.MustGet("token").(string)
		client := containersapi.NewContainersClient(token)

		serv, err := client.GetService(c, params.DBMS)
		if err != nil {
			return nil, err
		}

		inst, err := client.InstallService(c, serv.ID)
		if err != nil {
			return nil, err
		}

		inst.ContainerSettings.Tags = []string{"Vertex SQL", "Vertex SQL - Postgres Database"}
		err = client.PatchContainer(c, inst.UUID, inst.ContainerSettings)
		if err != nil {
			return nil, err
		}

		inst.Env, err = r.sqlService.EnvCredentials(inst, "postgres", "postgres")
		if err != nil {
			return nil, fmt.Errorf("setup credentials: %w", err)
		}

		err = client.PatchContainerEnvironment(c, inst.UUID, inst.Env)
		if err != nil {
			return nil, err
		}
		return inst, nil
	})
}
