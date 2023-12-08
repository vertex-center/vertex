package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/pkg/router"
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

type MigrateToParams struct {
	Dbms string `json:"dbms"`
}

func (h *databaseHandler) MigrateTo() gin.HandlerFunc {
	return router.Handler(func(c *gin.Context, params *MigrateToParams) error {
		return h.dataService.MigrateTo(params.Dbms)
	})
}
