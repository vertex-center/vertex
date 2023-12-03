package port

import (
	"context"

	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/app"
)

type (
	AppsService interface {
		All() []app.Meta
	}

	DebugService interface {
		HardReset()
	}

	ChecksService interface {
		CheckAll(ctx context.Context) <-chan types.CheckResponse
	}
)
