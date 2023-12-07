package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	ChecksHandler interface {
		Check(c *router.Context)
	}

	DatabaseHandler interface {
		GetCurrentDbms(c *router.Context)
		MigrateTo(c *router.Context)
	}

	HardwareHandler interface {
		GetHost(c *router.Context)
		GetCPUs(c *router.Context)
		Reboot(c *router.Context)
	}

	HardwareKernelHandler interface {
		Reboot(c *router.Context)
	}

	SettingsHandler interface {
		Get(c *router.Context)
		Patch(c *router.Context)
	}

	SshHandler interface {
		Get(c *router.Context)
		Add(c *router.Context)
		Delete(c *router.Context)
		GetUsers(c *router.Context)
	}

	SshKernelHandler interface {
		Get(c *router.Context)
		Add(c *router.Context)
		Delete(c *router.Context)
		GetUsers(c *router.Context)
	}

	UpdateHandler interface {
		Get(c *router.Context)
		Install(c *router.Context)
	}
)
