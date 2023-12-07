package handler

import (
	"net/http"

	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/core/types/api"
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

func (h *databaseHandler) GetCurrentDbms(c *router.Context) {
	c.JSON(h.dataService.GetCurrentDbms())
}

func (h *databaseHandler) GetCurrentDbmsInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Get the current DBMS"),
		oapi.Description("Get the current database management system that Vertex is using."),
		oapi.Response(http.StatusOK),
	}
}

type MigrateToBody struct {
	Dbms string `json:"dbms"`
}

func (h *databaseHandler) MigrateTo(c *router.Context) {
	var body MigrateToBody
	err := c.ParseBody(&body)
	if err != nil {
		return
	}

	err = h.dataService.MigrateTo(body.Dbms)
	//if err != nil && errors.Is(err, service.ErrDbmsAlreadySet) {
	//	c.NotModified()
	//	return
	//}
	if err != nil {
		c.Abort(router.Error{
			Code:           api.ErrFailedToMigrateToNewDbms,
			PublicMessage:  "Migration to the new DBMS failed.",
			PrivateMessage: err.Error(),
		})
		return
	}

	c.OK()
}

func (h *databaseHandler) MigrateToInfo() []oapi.Info {
	return []oapi.Info{
		oapi.Summary("Migrate to a DBMS"),
		oapi.Description("Migrate Vertex to the given database management system."),
		oapi.Response(http.StatusNoContent),
	}
}
