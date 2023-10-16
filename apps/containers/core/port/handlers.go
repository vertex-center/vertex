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
)
