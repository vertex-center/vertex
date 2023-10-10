package router

import (
	"errors"
	"fmt"

	instancesapi "github.com/vertex-center/vertex/apps/instances/api"
	instancestypes "github.com/vertex-center/vertex/apps/instances/types"
	"github.com/vertex-center/vertex/apps/sql/service"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
	"github.com/vertex-center/vertex/types/app"
)

type AppRouter struct {
	sqlService *service.SqlService
}

func NewAppRouter() *AppRouter {
	return &AppRouter{
		sqlService: service.New(),
	}
}

func (r *AppRouter) GetServices() []app.Service {
	return []app.Service{
		r.sqlService,
	}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.GET("/instance/:instance_uuid", r.handleGetDBMS)
	group.POST("/dbms/:dbms/install", r.handleInstallDBMS)
}

func (r *AppRouter) getDBMS(c *router.Context) (string, error) {
	db := c.Param("dbms")
	if db != "postgres" {
		c.NotFound(router.Error{
			Code:           api.ErrSQLDatabaseNotFound,
			PublicMessage:  fmt.Sprintf("SQL DBMS not found: %s.", db),
			PrivateMessage: "This SQL DBMS is not supported.",
		})
		return "", errors.New("DBMS not found")
	}

	return db, nil
}

func (r *AppRouter) handleGetDBMS(c *router.Context) {
	uuid, apiError := instancesapi.GetInstanceUUIDParam(c)
	if apiError != nil {
		c.BadRequest(apiError.RouterError())
		return
	}

	inst, apiError := instancesapi.GetInstance(c, uuid)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	db, err := r.sqlService.Get(inst)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrSQLDatabaseNotFound,
			PublicMessage:  "SQL Database not found.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.JSON(db)
}

func (r *AppRouter) handleInstallDBMS(c *router.Context) {
	db, err := r.getDBMS(c)
	if err != nil {
		return
	}

	serv, apiError := instancesapi.GetService(c, db)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	inst, apiError := instancesapi.InstallService(c, serv.ID)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	apiError = instancesapi.PatchInstance(c, inst.UUID, instancestypes.InstanceSettings{
		Tags: []string{"vertex-postgres-sql"},
	})
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	var env instancestypes.InstanceEnvVariables
	env, err = r.sqlService.EnvCredentials(inst, "postgres", "postgres")
	if err != nil {
		log.Error(err)
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureSQLDatabaseInstance,
			PublicMessage:  fmt.Sprintf("Failed to configure SQL Database '%s'.", serv.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	apiError = instancesapi.PatchInstanceEnvironment(c, inst.UUID, env)
	if apiError != nil {
		c.AbortWithCode(apiError.HttpCode, apiError.RouterError())
		return
	}

	c.OK()
}
