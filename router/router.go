package router

import (
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("vertex::router")

func InitializeRouter() *gin.Engine {
	r, api := router.CreateRouter(cors.Default())
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.PathClient, "dist"), true)))

	addServicesRoutes(api.Group("/services"))
	addInstancesRoutes(api.Group("/instances"))
	addInstanceRoutes(api.Group("/instance/:instance_uuid"))
	addDependenciesRoutes(api.Group("/dependencies"))

	return r
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
