package service

import (
	"github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
)

type debugService struct {
	ctx *types.VertexContext
}

func NewDebugService(ctx *types.VertexContext) port.DebugService {
	return &debugService{
		ctx: ctx,
	}
}

func (s *debugService) HardReset() {
	s.ctx.DispatchEvent(event.ServerHardReset{})
}
