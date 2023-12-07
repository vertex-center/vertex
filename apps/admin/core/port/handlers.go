package port

import (
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vertex/pkg/router/oapi"
)

type (
	ChecksHandler interface {
		Check(c *router.Context)
		CheckInfo() []oapi.Info
	}

	DatabaseHandler interface {
		GetCurrentDbms(c *router.Context)
		GetCurrentDbmsInfo() []oapi.Info

		MigrateTo(c *router.Context)
		MigrateToInfo() []oapi.Info
	}

	HardwareHandler interface {
		GetHost(c *router.Context)
		GetHostInfo() []oapi.Info

		GetCPUs(c *router.Context)
		GetCPUsInfo() []oapi.Info

		Reboot(c *router.Context)
		RebootInfo() []oapi.Info
	}

	HardwareKernelHandler interface {
		Reboot(c *router.Context)
		RebootInfo() []oapi.Info
	}

	SettingsHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Patch(c *router.Context)
		PatchInfo() []oapi.Info
	}

	SshHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Add(c *router.Context)
		AddInfo() []oapi.Info

		Delete(c *router.Context)
		DeleteInfo() []oapi.Info

		GetUsers(c *router.Context)
		GetUsersInfo() []oapi.Info
	}

	SshKernelHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Add(c *router.Context)
		AddInfo() []oapi.Info

		Delete(c *router.Context)
		DeleteInfo() []oapi.Info

		GetUsers(c *router.Context)
		GetUsersInfo() []oapi.Info
	}

	UpdateHandler interface {
		Get(c *router.Context)
		GetInfo() []oapi.Info

		Install(c *router.Context)
		InstallInfo() []oapi.Info
	}
)
