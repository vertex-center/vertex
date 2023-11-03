package app

import (
	"github.com/vertex-center/vertex/core/types"
	evtypes "github.com/vertex-center/vertex/pkg/event/types"
)

type Context struct {
	vertexCtx *types.VertexContext
}

func NewContext(vertexCtx *types.VertexContext) *Context {
	return &Context{
		vertexCtx: vertexCtx,
	}
}

func (ctx *Context) AddListener(listener evtypes.EventListener) {
	ctx.vertexCtx.AddListener(listener)
}

func (ctx *Context) RemoveListener(listener evtypes.EventListener) {
	ctx.vertexCtx.RemoveListener(listener)
}

func (ctx *Context) DispatchEvent(event evtypes.Event) {
	ctx.vertexCtx.DispatchEvent(event)
}
