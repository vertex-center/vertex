package app

import (
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/event"
)

type Context struct {
	vertexCtx *types.VertexContext
}

func NewContext(vertexCtx *types.VertexContext) *Context {
	return &Context{
		vertexCtx: vertexCtx,
	}
}

func (ctx *Context) AddListener(listener event.Listener) {
	ctx.vertexCtx.AddListener(listener)
}

func (ctx *Context) RemoveListener(listener event.Listener) {
	ctx.vertexCtx.RemoveListener(listener)
}

func (ctx *Context) DispatchEvent(event interface{}) {
	ctx.vertexCtx.DispatchEvent(event)
}
