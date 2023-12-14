package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ContainerHandler interface {
		Get() gin.HandlerFunc
		Delete() gin.HandlerFunc
		Patch() gin.HandlerFunc
		Start() gin.HandlerFunc
		Stop() gin.HandlerFunc
		AddTag() gin.HandlerFunc
		PatchEnvironment() gin.HandlerFunc
		GetDocker() gin.HandlerFunc
		RecreateDocker() gin.HandlerFunc
		GetLogs() gin.HandlerFunc
		GetVersions() gin.HandlerFunc
		WaitStatus() gin.HandlerFunc
		Events() gin.HandlerFunc
	}

	ContainersHandler interface {
		Get() gin.HandlerFunc
		Search() gin.HandlerFunc
		CheckForUpdates() gin.HandlerFunc
		Events() gin.HandlerFunc
	}

	ServiceHandler interface {
		Get() gin.HandlerFunc
		Install() gin.HandlerFunc
	}

	ServicesHandler interface {
		Get() gin.HandlerFunc
	}

	TagsHandler interface {
		GetTags() gin.HandlerFunc
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
