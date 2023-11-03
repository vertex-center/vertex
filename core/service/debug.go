package service

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
)

type DebugService struct {
	ctx *types.VertexContext
}

func NewDebugService(ctx *types.VertexContext) port.DebugService {
	return &DebugService{
		ctx: ctx,
	}
}

func (s *DebugService) HardReset() {
	s.ctx.DispatchEvent(types.EventServerHardReset{})
}
