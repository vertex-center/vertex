package router

import (
	"errors"
	"fmt"

	"github.com/carlmjohnson/requests"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/sql/service"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vertex/types/api"
)

type AppRouter struct {
	sqlService *service.SqlService
}

func NewAppRouter() *AppRouter {
	return &AppRouter{
		sqlService: service.NewSqlService(),
	}
}

func (r *AppRouter) GetServices() []types.AppService {
	return []types.AppService{
		r.sqlService,
	}
}

func (r *AppRouter) AddRoutes(group *router.Group) {
	group.GET("/:instance_uuid", r.handleGetDBMS)
	group.POST("/db/:db/install", r.handleInstallDBMS)
}

func (r *AppRouter) getDBMS(c *router.Context) (string, error) {
	db := c.Param("db")
	if db != "postgres" {
		c.NotFound(router.Error{
			Code:           api.ErrSQLDatabaseNotFound,
			PublicMessage:  fmt.Sprintf("SQL Database not found: %s.", db),
			PrivateMessage: "This SQL Database is not supported.",
		})
		return "", errors.New("collector not found")
	}

	return db, nil
}

func (r *AppRouter) handleGetDBMS(c *router.Context) {
	uuid := c.Param("instance_uuid")
	if uuid == "" {
		c.BadRequest(router.Error{
			Code:           api.ErrInstanceUuidMissing,
			PublicMessage:  "The request was missing the instance UUID.",
			PrivateMessage: "Field 'instance_uuid' is required.",
		})
		return
	}

	var inst *types.Instance
	var apiError router.Error
	err := requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s", uuid).
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(c)
	if err != nil {
		log.Error(err)
		c.NotFound(apiError)
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

	var serv types.Service
	var apiError router.Error
	err = requests.URL(config.Current.VertexURL()).
		Pathf("/api/service/%s", db).
		ToJSON(&serv).
		ErrorJSON(&apiError).
		Fetch(c)
	if err != nil {
		log.Error(err)
		c.NotFound(apiError)
		return
	}

	var inst *types.Instance
	err = requests.URL(config.Current.VertexURL()).
		Pathf("/api/service/%s/install", db).
		Post().
		ToJSON(&inst).
		ErrorJSON(&apiError).
		Fetch(c)
	if err != nil {
		log.Error(err)
		c.NotFound(apiError)
		return
	}

	err = requests.URL(config.Current.VertexURL()).
		Pathf("/api/instance/%s", inst.UUID).
		Patch().
		BodyJSON(gin.H{
			"tags": []string{"vertex-postgres-sql"},
		}).
		ErrorJSON(&apiError).
		Fetch(c)

	var env types.InstanceEnvVariables
	if err == nil {
		env, err = r.sqlService.EnvCredentials(inst, "postgres", "postgres")
	}
	if err == nil {
		err = requests.URL(config.Current.VertexURL()).
			Pathf("/api/instance/%s/environment", inst.UUID).
			Patch().
			BodyJSON(&env).
			ErrorJSON(&apiError).
			Fetch(c)
	}

	if err != nil {
		log.Error(err)
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureSQLDatabaseInstance,
			PublicMessage:  fmt.Sprintf("Failed to configure SQL Database '%s'.", serv.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
