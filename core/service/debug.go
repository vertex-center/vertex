package service

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/log"
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
	err := s.ctx.DispatchEvent(types.EventServerHardReset{})
	if err != nil {
		log.Error(err)
	}
}
