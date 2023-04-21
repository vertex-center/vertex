package router

import (
	"net/http"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/console"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/services"
	"github.com/vertex-center/vertex/storage"
)

var logger = console.New("vertex::router")

var (
	packageService services.PackageService
	serviceService services.ServiceService
)

type About struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func InitializeRouter(about About) *gin.Engine {
	r, api := router.CreateRouter(cors.Default())
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.PathClient, "dist"), true)))

	packageService = services.NewPackageService()
	serviceService = services.NewServiceService()

	addServicesRoutes(api.Group("/services"))
	addInstancesRoutes(api.Group("/instances"))
	addInstanceRoutes(api.Group("/instance/:instance_uuid"))
	addPackagesRoutes(api.Group("/packages"))
	addUpdatesRoutes(api.Group("/updates"), about.Version)

	api.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK, about)
	})

	return r
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
