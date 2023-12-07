package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type databaseHandler struct {
	dataService port.DatabaseService
}

func NewDatabaseHandler(dataService port.DatabaseService) port.DatabaseHandler {
	return &databaseHandler{
		dataService: dataService,
	}
}

func (h *databaseHandler) GetCurrentDbms() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context) (*string, error) {
		dbms := h.dataService.GetCurrentDbms()
		return &dbms, nil
	})
}

func (h *databaseHandler) GetCurrentDbmsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("getCurrentDbms"),
		oapi.Summary("Get the current DBMS"),
		oapi.Description("Get the current database management system that Vertex is using."),
	}
}

type MigrateToParams struct {
	Dbms string `json:"dbms"`
}

func (h *databaseHandler) MigrateTo() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *MigrateToParams) error {
		return h.dataService.MigrateTo(params.Dbms)
	})
}

func (h *databaseHandler) MigrateToInfo() []oapi.Info {
	return []oapi.Info{
		oapi.ID("migrateTo"),
		oapi.Summary("Migrate to a DBMS"),
		oapi.Description("Migrate Vertex to the given database management system."),
	}
}
