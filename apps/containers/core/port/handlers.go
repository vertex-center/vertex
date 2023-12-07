package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ContainerHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Delete(c *router.Context)
		DeleteInfo() []oapi.Info

		Patch(c *router.Context)
		PatchInfo() []oapi.Info

		Start(c *router.Context)
		StartInfo() []oapi.Info

		Stop(c *router.Context)
		StopInfo() []oapi.Info

		PatchEnvironment(c *router.Context)
		PatchEnvironmentInfo() []oapi.Info

		GetDocker(c *router.Context)
		GetDockerInfo() []oapi.Info

		RecreateDocker(c *router.Context)
		RecreateDockerInfo() []oapi.Info

		GetLogs(c *router.Context)
		GetLogsInfo() []oapi.Info

		UpdateService(c *router.Context)
		UpdateServiceInfo() []oapi.Info

		GetVersions(c *router.Context)
		GetVersionsInfo() []oapi.Info

		WaitStatus(c *router.Context)
		WaitStatusInfo() []oapi.Info

		Events(c *router.Context)
		EventsInfo() []oapi.Info
	}

	ContainersHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		GetTags(c *router.Context)
		GetTagsInfo() []oapi.Info

		Search(c *router.Context)
		SearchInfo() []oapi.Info

		CheckForUpdates(c *router.Context)
		CheckForUpdatesInfo() []oapi.Info

		Events(c *router.Context)
		EventsInfo() []oapi.Info
	}

	ServiceHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Install(c *router.Context)
		InstallInfo() []oapi.Info
	}

	ServicesHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info
	}

	DockerKernelHandler interface {
		GetContainers(c *router.Context)
		GetContainersInfo() []oapi.Info

		CreateContainer(c *router.Context)
		CreateContainerInfo() []oapi.Info

		DeleteContainer(c *router.Context)
		DeleteContainerInfo() []oapi.Info

		StartContainer(c *router.Context)
		StartContainerInfo() []oapi.Info

		StopContainer(c *router.Context)
		StopContainerInfo() []oapi.Info

		InfoContainer(c *router.Context)
		InfoContainerInfo() []oapi.Info

		LogsStdoutContainer(c *router.Context)
		LogsStdoutContainerInfo() []oapi.Info

		LogsStderrContainer(c *router.Context)
		LogsStderrContainerInfo() []oapi.Info

		WaitContainer(c *router.Context)
		WaitContainerInfo() []oapi.Info

		DeleteMounts(c *router.Context)
		DeleteMountsInfo() []oapi.Info

		InfoImage(c *router.Context)
		InfoImageInfo() []oapi.Info

		PullImage(c *router.Context)
		PullImageInfo() []oapi.Info

		BuildImage(c *router.Context)
		BuildImageInfo() []oapi.Info
	}
)
