package router

import (
	"errors"
	"fmt"

	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/types/api"
)

func addSQLRoutes(r *router.Group) {
	r.POST("/db/:db/install", handleInstallSQLDatabase)
}

func getSQLDatabase(c *router.Context) (string, error) {
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

func handleInstallSQLDatabase(c *router.Context) {
	db, err := getSQLDatabase(c)
	if err != nil {
		return
	}

	service, err := serviceService.GetById(db)
	if err != nil {
		c.NotFound(router.Error{
			Code:           api.ErrServiceNotFound,
			PublicMessage:  fmt.Sprintf("Service not found: %s.", db),
			PrivateMessage: err.Error(),
		})
		return
	}

	inst, err := instanceService.Install(service, "docker")
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToInstallService,
			PublicMessage:  fmt.Sprintf("Failed to install service '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	err = instanceSettingsService.SetTags(inst, []string{"vertex-postgres-sql"})
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToConfigureSQLDatabaseInstance,
			PublicMessage:  fmt.Sprintf("Failed to configure SQL Database '%s'.", service.Name),
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}
