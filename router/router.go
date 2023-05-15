package router

import (
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sse"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex-core-golang/router"
	"github.com/vertex-center/vertex/pkg/logger"
	"github.com/vertex-center/vertex/pkg/storage"
	"github.com/vertex-center/vertex/repository"
	"github.com/vertex-center/vertex/services"
)

var (
	runnerDockerRepo  repository.RunnerDockerRepository
	runnerFSRepo      repository.RunnerFSRepository
	instanceRepo      repository.InstanceFSRepository
	instanceLogsRepo  repository.InstanceLogsFSRepository
	eventInMemoryRepo repository.EventInMemoryRepository
	packageRepo       repository.PackageFSRepository
	serviceRepo       repository.ServiceFSRepository

	packageService  services.PackageService
	serviceService  services.ServiceService
	instanceService services.InstanceService
	updateService   services.UpdateDependenciesService
)

type About struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

func Create(about About) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r, api := router.CreateRouter(
		cors.Default(),
		gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
			l := logger.Request().
				AddKeyValue("method", params.Method).
				AddKeyValue("status", params.StatusCode).
				AddKeyValue("path", params.Path).
				AddKeyValue("latency", params.Latency).
				AddKeyValue("ip", params.ClientIP).
				AddKeyValue("size", params.BodySize)

			if params.ErrorMessage != "" {
				err, _ := strings.CutSuffix(params.ErrorMessage, "\n")
				l.AddKeyValue("error", err)
			}

			l.PrintInExternalFiles()

			return l.String()
		}),
	)
	r.Use(static.Serve("/", static.LocalFile(path.Join(".", storage.PathClient, "dist"), true)))

	runnerDockerRepo = repository.NewRunnerDockerRepository()
	runnerFSRepo = repository.NewRunnerFSRepository()
	instanceRepo = repository.NewInstanceFSRepository()
	instanceLogsRepo = repository.NewInstanceLogsFSRepository()
	eventInMemoryRepo = repository.NewEventInMemoryRepository()
	packageRepo = repository.NewPackageFSRepository(nil)
	serviceRepo = repository.NewServiceFSRepository(nil)

	instanceService = services.NewInstanceService(&instanceRepo, &runnerDockerRepo, &runnerFSRepo, &instanceLogsRepo, &eventInMemoryRepo)
	packageService = services.NewPackageService(&packageRepo)
	serviceService = services.NewServiceService(&serviceRepo)
	updateService = services.NewUpdateDependenciesService(about.Version)

	go func() {
		instanceService.StartAll()
	}()

	handleSignals()

	addServicesRoutes(api.Group("/services"))
	addInstancesRoutes(api.Group("/instances"))
	addInstanceRoutes(api.Group("/instance/:instance_uuid"))
	addPackagesRoutes(api.Group("/packages"))
	addUpdatesRoutes(api.Group("/updates"))

	api.GET("/about", func(c *gin.Context) {
		c.JSON(http.StatusOK, about)
	})

	return r
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		logger.Log("shutdown signal sent").Print()
		Unload()
		os.Exit(0)
	}()
}

func Unload() {
	instanceService.StopAll()
}

func headersSSE(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
