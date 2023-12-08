package port

import (
	"github.com/gin-gonic/gin"
	"github.com/wI2L/fizz"
)

type (
	ContainerHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Delete() gin.HandlerFunc
		DeleteInfo() []fizz.OperationOption

		Patch() gin.HandlerFunc
		PatchInfo() []fizz.OperationOption

		Start() gin.HandlerFunc
		StartInfo() []fizz.OperationOption

		Stop() gin.HandlerFunc
		StopInfo() []fizz.OperationOption

		PatchEnvironment() gin.HandlerFunc
		PatchEnvironmentInfo() []fizz.OperationOption

		GetDocker() gin.HandlerFunc
		GetDockerInfo() []fizz.OperationOption

		RecreateDocker() gin.HandlerFunc
		RecreateDockerInfo() []fizz.OperationOption

		GetLogs() gin.HandlerFunc
		GetLogsInfo() []fizz.OperationOption

		UpdateService() gin.HandlerFunc
		UpdateServiceInfo() []fizz.OperationOption

		GetVersions() gin.HandlerFunc
		GetVersionsInfo() []fizz.OperationOption

		WaitStatus() gin.HandlerFunc
		WaitStatusInfo() []fizz.OperationOption

		Events() gin.HandlerFunc
		EventsInfo() []fizz.OperationOption
	}

	ContainersHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		GetTags() gin.HandlerFunc
		GetTagsInfo() []fizz.OperationOption

		Search() gin.HandlerFunc
		SearchInfo() []fizz.OperationOption

		CheckForUpdates() gin.HandlerFunc
		CheckForUpdatesInfo() []fizz.OperationOption

		Events() gin.HandlerFunc
		EventsInfo() []fizz.OperationOption
	}

	ServiceHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption

		Install() gin.HandlerFunc
		InstallInfo() []fizz.OperationOption
	}

	ServicesHandler interface {
		Get() gin.HandlerFunc
		GetInfo() []fizz.OperationOption
	}

	DockerKernelHandler interface {
		GetContainers() gin.HandlerFunc
		GetContainersInfo() []fizz.OperationOption

		CreateContainer() gin.HandlerFunc
		CreateContainerInfo() []fizz.OperationOption

		DeleteContainer() gin.HandlerFunc
		DeleteContainerInfo() []fizz.OperationOption

		StartContainer() gin.HandlerFunc
		StartContainerInfo() []fizz.OperationOption

		StopContainer() gin.HandlerFunc
		StopContainerInfo() []fizz.OperationOption

		InfoContainer() gin.HandlerFunc
		InfoContainerInfo() []fizz.OperationOption

		LogsStdoutContainer() gin.HandlerFunc
		LogsStdoutContainerInfo() []fizz.OperationOption

		LogsStderrContainer() gin.HandlerFunc
		LogsStderrContainerInfo() []fizz.OperationOption

		WaitContainer() gin.HandlerFunc
		WaitContainerInfo() []fizz.OperationOption

		DeleteMounts() gin.HandlerFunc
		DeleteMountsInfo() []fizz.OperationOption

		InfoImage() gin.HandlerFunc
		InfoImageInfo() []fizz.OperationOption

		PullImage() gin.HandlerFunc
		PullImageInfo() []fizz.OperationOption

		BuildImage() gin.HandlerFunc
		BuildImageInfo() []fizz.OperationOption
	}
)
