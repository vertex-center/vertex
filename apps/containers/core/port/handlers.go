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
		Wait(c *router.Context)
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
		// GetContainers handles the retrieval of all Docker containers.
		GetContainers(c *router.Context)
		// CreateContainer handles the creation of a Docker container.
		CreateContainer(c *router.Context)
		// DeleteContainer handles the deletion of a Docker container.
		DeleteContainer(c *router.Context)
		// StartContainer handles the starting of a Docker container.
		StartContainer(c *router.Context)
		// StopContainer handles the stopping of a Docker container.
		StopContainer(c *router.Context)
		// InfoContainer handles the retrieval of information about a Docker container.
		InfoContainer(c *router.Context)
		// LogsStdoutContainer handles the retrieval of the stdout logs of a Docker container.
		LogsStdoutContainer(c *router.Context)
		// LogsStderrContainer handles the retrieval of the stderr logs of a Docker container.
		LogsStderrContainer(c *router.Context)
		// WaitContainer handles the waiting for a Docker container to reach a certain condition.
		WaitContainer(c *router.Context)
		// InfoImage handles the retrieval of information about a Docker image.
		InfoImage(c *router.Context)
		// PullImage handles the pulling of a Docker image.
		PullImage(c *router.Context)
		// BuildImage handles the building of a Docker image.
		BuildImage(c *router.Context)
	}
)
