package port

import (
	"github.com/gin-gonic/gin"
)

type (
	ContainerHandler interface {
		Get() gin.HandlerFunc
		GetContainers() gin.HandlerFunc
		CreateContainer() gin.HandlerFunc
		Delete() gin.HandlerFunc
		Patch() gin.HandlerFunc
		Start() gin.HandlerFunc
		Stop() gin.HandlerFunc
		AddContainerTag() gin.HandlerFunc
		GetDocker() gin.HandlerFunc
		RecreateDocker() gin.HandlerFunc
		ReloadContainer() gin.HandlerFunc
		GetLogs() gin.HandlerFunc
		GetVersions() gin.HandlerFunc
		WaitStatus() gin.HandlerFunc
		CheckForUpdates() gin.HandlerFunc
		ContainerEvents() gin.HandlerFunc
		ContainersEvents() gin.HandlerFunc
	}

	EnvHandler interface {
		GetEnv() gin.HandlerFunc
		PatchEnv() gin.HandlerFunc
		DeleteEnv() gin.HandlerFunc
		CreateEnv() gin.HandlerFunc
	}

	PortsHandler interface {
		GetPorts() gin.HandlerFunc
		PatchPort() gin.HandlerFunc
		DeletePort() gin.HandlerFunc
		CreatePort() gin.HandlerFunc
	}

	TemplateHandler interface {
		GetTemplate() gin.HandlerFunc
		GetTemplates() gin.HandlerFunc
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
		CreateVolume() gin.HandlerFunc
		DeleteVolume() gin.HandlerFunc
	}
)
