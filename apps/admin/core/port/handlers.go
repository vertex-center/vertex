package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	DatabaseHandler interface {
		// GetCurrentDbms handles the retrieval of the current database management system
		// that Vertex is using.
		GetCurrentDbms(c *router.Context)
		// MigrateTo handles the migration to the given database management system.
		MigrateTo(c *router.Context)
	}

	HardwareHandler interface {
		// GetHost handles the retrieval of the current host.
		GetHost(c *router.Context)
		// GetCPUs handles the retrieval of all CPUs.
		GetCPUs(c *router.Context)
		// Reboot handles the reboot of the server.
		Reboot(c *router.Context)
	}

	HardwareKernelHandler interface {
		// Reboot handles the reboot of the server.
		Reboot(c *router.Context)
	}

	SettingsHandler interface {
		// Get handles the retrieval of all settings.
		Get(c *router.Context)
		// Patch handles the update of all settings.
		Patch(c *router.Context)
	}

	SshHandler interface {
		// Get handles the retrieval of all SSH keys.
		Get(c *router.Context)
		// Add handles the addition of an SSH key.
		Add(c *router.Context)
		// Delete handles the deletion of an SSH key.
		Delete(c *router.Context)
		// GetUsers handles the retrieval of all users that have can have SSH keys.
		GetUsers(c *router.Context)
	}

	SshKernelHandler interface {
		// Get handles the retrieval of all SSH keys.
		Get(c *router.Context)
		// Add handles the addition of an SSH key.
		Add(c *router.Context)
		// Delete handles the deletion of an SSH key.
		Delete(c *router.Context)
		// GetUsers handles the retrieval of all users that have can have SSH keys.
		GetUsers(c *router.Context)
	}

	UpdateHandler interface {
		// Get handles the retrieval of an update, if any.
		Get(c *router.Context)
		// Install handles the installation of the update.
		Install(c *router.Context)
	}
)
