package port

import (
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ContainerHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Delete() gin.HandlerFunc
		DeleteInfo() []oapi.Info

		Patch() gin.HandlerFunc
		PatchInfo() []oapi.Info

		Start() gin.HandlerFunc
		StartInfo() []oapi.Info

		Stop() gin.HandlerFunc
		StopInfo() []oapi.Info

		PatchEnvironment() gin.HandlerFunc
		PatchEnvironmentInfo() []oapi.Info

		GetDocker() gin.HandlerFunc
		GetDockerInfo() []oapi.Info

		RecreateDocker() gin.HandlerFunc
		RecreateDockerInfo() []oapi.Info

		GetLogs() gin.HandlerFunc
		GetLogsInfo() []oapi.Info

		UpdateService() gin.HandlerFunc
		UpdateServiceInfo() []oapi.Info

		GetVersions() gin.HandlerFunc
		GetVersionsInfo() []oapi.Info

		WaitStatus() gin.HandlerFunc
		WaitStatusInfo() []oapi.Info

		Events() gin.HandlerFunc
		EventsInfo() []oapi.Info
	}

	ContainersHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		GetTags() gin.HandlerFunc
		GetTagsInfo() []oapi.Info

		Search() gin.HandlerFunc
		SearchInfo() []oapi.Info

		CheckForUpdates() gin.HandlerFunc
		CheckForUpdatesInfo() []oapi.Info

		Events() gin.HandlerFunc
		EventsInfo() []oapi.Info
	}

	ServiceHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info

		Install() gin.HandlerFunc
		InstallInfo() []oapi.Info
	}

	ServicesHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []oapi.Info
	}

	DockerKernelHandler interface {
		GetContainers() gin.HandlerFunc
		GetContainersInfo() []oapi.Info

		CreateContainer() gin.HandlerFunc
		CreateContainerInfo() []oapi.Info

		DeleteContainer() gin.HandlerFunc
		DeleteContainerInfo() []oapi.Info

		StartContainer() gin.HandlerFunc
		StartContainerInfo() []oapi.Info

		StopContainer() gin.HandlerFunc
		StopContainerInfo() []oapi.Info

		InfoContainer() gin.HandlerFunc
		InfoContainerInfo() []oapi.Info

		LogsStdoutContainer() gin.HandlerFunc
		LogsStdoutContainerInfo() []oapi.Info

		LogsStderrContainer() gin.HandlerFunc
		LogsStderrContainerInfo() []oapi.Info

		WaitContainer() gin.HandlerFunc
		WaitContainerInfo() []oapi.Info

		DeleteMounts() gin.HandlerFunc
		DeleteMountsInfo() []oapi.Info

		InfoImage() gin.HandlerFunc
		InfoImageInfo() []oapi.Info

		PullImage() gin.HandlerFunc
		PullImageInfo() []oapi.Info

		BuildImage() gin.HandlerFunc
		BuildImageInfo() []oapi.Info
	}
)
