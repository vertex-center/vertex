package handler

import (
	"errors"
	"fmt"
	"net/http"

	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	containerstypes "github.com/vertex-center/vertex/apps/containers/core/types"
	"github.com/vertex-center/vertex/apps/sql/core/port"
	"github.com/vertex-center/vertex/apps/sql/core/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type dbmsHandler struct {
	sqlService port.SqlService
}

func NewDBMSHandler(sqlService port.SqlService) port.DBMSHandler {
	return &dbmsHandler{
		sqlService: sqlService,
	}
}

func (r *dbmsHandler) Get(c *router.Context) {
	uuid, apiError := containersapi.GetContainerUUIDParam(c)
	if apiError != nil {
		c.BadRequest(apiError.RouterError())
		return
	}

	token := c.MustGet("token").(string)

	client := containersapi.NewContainersClient(token)

	inst, apiError := client.GetContainer(c, uuid)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	dbms, err := r.sqlService.Get(inst)
	if err != nil {
		c.NotFound(router.Error{
			Code:           types.ErrCodeSQLDatabaseNotFound,
			PublicMessage:  "SQL Database not found.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(dbms)
}

func (r *dbmsHandler) GetInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get an installed DBMS"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(types.DBMS{}),
		),
	}
}

func (r *dbmsHandler) Install(c *router.Context) {
	dbms, err := r.getDBMS(c)
	if err != nil {
		return
	}

	token := c.MustGet("token").(string)

	client := containersapi.NewContainersClient(token)

	serv, apiError := client.GetService(c, dbms)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := client.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst.ContainerSettings.Tags = []string{"Vertex SQL", "Vertex SQL - Postgres Database"}
	apiError = client.PatchContainer(c, inst.UUID, inst.ContainerSettings)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst.Env, err = r.sqlService.EnvCredentials(inst, "postgres", "postgres")
	if err != nil {
		log.Error(err)
		c.Abort(router.Error{
			Code:           types.ErrCodeFailedToConfigureSQLDatabaseContainer,
			PublicMessage:  fmt.Sprintf("Failed to configure SQL Database '%s'.", serv.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	apiError = client.PatchContainerEnvironment(c, inst.UUID, inst.Env)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.JSON(inst)
}

func (r *dbmsHandler) InstallInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Install a DBMS"),
		oapi.Response(http.StatusOK,
			oapi.WithResponseModel(containerstypes.Container{}),
		),
	}
}

func (r *dbmsHandler) getDBMS(c *router.Context) (string, error) {
	db := c.Param("dbms")
	if db != "postgres" {
		c.NotFound(router.Error{
			Code:           types.ErrCodeSQLDatabaseNotFound,
			PublicMessage:  fmt.Sprintf("SQL DBMS not found: %s.", db),
			PrivateMessage: "This SQL DBMS is not supported.",
		})
		return "", errors.New("DBMS not found")
	}

	return db, nil
}
