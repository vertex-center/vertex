package port

import "github.com/vertex-center/vertex/pkg/router"

type AppsHandler interface {
	// Get handles the retrieval of all apps.
	Get(c *router.Context)
}

type HardwareHandler interface {
	// Get handles the retrieval of the current hardware.
	Get(c *router.Context)
}

type UpdateHandler interface {
	// Get handles the retrieval of an update, if any.
	Get(c *router.Context)
	// Install handles the installation of the update.
	Install(c *router.Context)
}

type SettingsHandler interface {
	// Get handles the retrieval of all settings.
	Get(c *router.Context)
	// Patch handles the update of all settings.
	Patch(c *router.Context)
}

type SshHandler interface {
	// Get handles the retrieval of all SSH keys.
	Get(c *router.Context)
	// Add handles the addition of an SSH key.
	Add(c *router.Context)
	// Delete handles the deletion of an SSH key.
	Delete(c *router.Context)
}

type DockerKernelHandler interface {
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

type SshKernelHandler interface {
	// Get handles the retrieval of all SSH keys.
	Get(c *router.Context)
	// Add handles the addition of an SSH key.
	Add(c *router.Context)
	// Delete handles the deletion of an SSH key.
	Delete(c *router.Context)
}
