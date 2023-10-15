package router

import (
	"errors"
	"fmt"
	containersapi "github.com/vertex-center/vertex/apps/containers/api"
	"github.com/vertex-center/vertex/apps/sql/service"
	"github.com/vertex-center/vertex/apps/sql/types"
	app2 "github.com/vertex-center/vertex/core/types/app"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
)

type AppRouter struct {
	sqlService *service.SqlService
}

func NewAppRouter(ctx *app2.Context) *AppRouter {
	return &AppRouter{
		sqlService: service.New(ctx),
	}
}

func (r *AppRouter) GetServices() []app2.Service {
	return []app2.Service{
		r.sqlService,
	}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.GET("/container/:container_uuid", r.handleGetDBMS)
	group.POST("/dbms/:dbms/install", r.handleInstallDBMS)
}

func (r *AppRouter) getDBMS(c *router.Context) (string, error) {
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

func (r *AppRouter) handleGetDBMS(c *router.Context) {
	uuid, apiError := containersapi.GetContainerUUIDParam(c)
	if apiError != nil {
		c.BadRequest(apiError.RouterError())
		return
	}

	inst, apiError := containersapi.GetContainer(c, uuid)
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

func (r *AppRouter) handleInstallDBMS(c *router.Context) {
	dbms, err := r.getDBMS(c)
	if err != nil {
		return
	}

	serv, apiError := containersapi.GetService(c, dbms)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := containersapi.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst.ContainerSettings.Tags = []string{"Vertex SQL", "Vertex SQL - Postgres Database"}
	apiError = containersapi.PatchContainer(c, inst.UUID, inst.ContainerSettings)
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

	apiError = containersapi.PatchContainerEnvironment(c, inst.UUID, inst.Env)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.JSON(inst)
}
