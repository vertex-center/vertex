package port

import (
	"github.com/vertex-center/vertex/apps/admin/core/types"
)

type (
	HardwareService interface {
		GetHost() (types.Host, error)
		GetCPUs() ([]types.CPU, error)
	}
)
