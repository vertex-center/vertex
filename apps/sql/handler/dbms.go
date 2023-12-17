package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/common/uuid"
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
	UUID uuid.UUID `path:"container_uuid"`
}

func (h *dbmsHandler) Get() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *GetParams) (*types.DBMS, error) {
		client := containersapi.NewContainersClient(c)

		inst, err := client.GetContainer(c, params.UUID)
		if err != nil {
			return nil, err
		}

		dbms, err := h.sqlService.Get(inst)
		if err != nil {
			return nil, errors.NewNotFound(err, "SQL Database not found")
		}
		return &dbms, nil
	})
}

type InstallParams struct {
	DBMS string `path:"dbms"`
}

func (h *dbmsHandler) Install() gin.HandlerFunc {
	return router.Handler(func(ctx *gin.Context, params *InstallParams) (*containerstypes.Container, error) {
		c, err := h.sqlService.Install(ctx, params.DBMS)
		if err != nil {
			return nil, err
		}
		return &c, nil
	})
}
