package port

import "github.com/vertex-center/vertex/pkg/router"

type (
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
