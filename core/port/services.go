package port

import (
	"github.com/vertex-center/vertex/core/types/app"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DebugService interface {
		HardReset()
	}
)
