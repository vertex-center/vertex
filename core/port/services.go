package port

import (
	"github.com/vertex-center/vertex/common/app"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DebugService interface {
		HardReset()
	}
)
