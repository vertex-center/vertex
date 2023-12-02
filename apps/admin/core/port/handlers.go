package port

import "github.com/vertex-center/vertex/pkg/router"

type (
	HardwareHandler interface {
		// GetHost handles the retrieval of the current host.
		GetHost(c *router.Context)
		// GetCPUs handles the retrieval of all CPUs.
		GetCPUs(c *router.Context)
	}
)
