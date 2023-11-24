package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	AppsHandler interface {
		// Get handles the retrieval of all apps.
		Get(c *router.Context)
	}

	DebugHandler interface {
		// HardReset do a hard reset of Vertex.
		HardReset(c *router.Context)
	}

	HardwareHandler interface {
		// Get handles the retrieval of the current hardware.
		Get(c *router.Context)
	}

	UpdateHandler interface {
		// Get handles the retrieval of an update, if any.
		Get(c *router.Context)
		// Install handles the installation of the update.
		Install(c *router.Context)
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
)
