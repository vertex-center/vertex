package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	ContainerHandler interface {
		Get(c *router.Context)
		Delete(c *router.Context)
		Patch(c *router.Context)
		Start(c *router.Context)
		Stop(c *router.Context)
		PatchEnvironment(c *router.Context)
		GetDocker(c *router.Context)
		RecreateDocker(c *router.Context)
		GetLogs(c *router.Context)
		UpdateService(c *router.Context)
		GetVersions(c *router.Context)
		WaitStatus(c *router.Context)
		Events(c *router.Context)
	}

	ContainersHandler interface {
		Get(c *router.Context)
		GetTags(c *router.Context)
		Search(c *router.Context)
		CheckForUpdates(c *router.Context)
		Events(c *router.Context)
	}

	ServiceHandler interface {
		Get(c *router.Context)
		Install(c *router.Context)
	}

	ServicesHandler interface {
		Get(c *router.Context)
	}

	DockerKernelHandler interface {
		GetContainers(c *router.Context)
		CreateContainer(c *router.Context)
		DeleteContainer(c *router.Context)
		StartContainer(c *router.Context)
		StopContainer(c *router.Context)
		InfoContainer(c *router.Context)
		LogsStdoutContainer(c *router.Context)
		LogsStderrContainer(c *router.Context)
		WaitContainer(c *router.Context)
		DeleteMounts(c *router.Context)

		InfoImage(c *router.Context)
		PullImage(c *router.Context)
		BuildImage(c *router.Context)
	}
)
