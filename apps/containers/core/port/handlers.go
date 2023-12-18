package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ContainerHandler interface {
		Get() gin.HandlerFunc
		CreateContainer() gin.HandlerFunc
		Delete() gin.HandlerFunc
		Patch() gin.HandlerFunc
		Start() gin.HandlerFunc
		Stop() gin.HandlerFunc
		AddContainerTag() gin.HandlerFunc
		GetContainerEnv() gin.HandlerFunc
		PatchEnvironment() gin.HandlerFunc
		GetDocker() gin.HandlerFunc
		RecreateDocker() gin.HandlerFunc
		GetLogs() gin.HandlerFunc
		GetVersions() gin.HandlerFunc
		WaitStatus() gin.HandlerFunc
		Events() gin.HandlerFunc
	}

	ContainersHandler interface {
		GetContainers() gin.HandlerFunc
		CheckForUpdates() gin.HandlerFunc
		Events() gin.HandlerFunc
	}

	ServiceHandler interface {
		Get() gin.HandlerFunc
	}

	ServicesHandler interface {
		Get() gin.HandlerFunc
	}

	TagsHandler interface {
		GetTag() gin.HandlerFunc
		GetTags() gin.HandlerFunc
		CreateTag() gin.HandlerFunc
		DeleteTag() gin.HandlerFunc
	}

	DockerKernelHandler interface {
		GetContainers() gin.HandlerFunc
		CreateContainer() gin.HandlerFunc
		DeleteContainer() gin.HandlerFunc
		StartContainer() gin.HandlerFunc
		StopContainer() gin.HandlerFunc
		InfoContainer() gin.HandlerFunc
		LogsStdoutContainer() gin.HandlerFunc
		LogsStderrContainer() gin.HandlerFunc
		WaitContainer() gin.HandlerFunc
		DeleteMounts() gin.HandlerFunc
		InfoImage() gin.HandlerFunc
		PullImage() gin.HandlerFunc
		BuildImage() gin.HandlerFunc
	}
)
