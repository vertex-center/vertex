package app

import (
	types2 "github.com/vertex-center/vertex/core/types"
)

type Context struct {
	vertexCtx *types2.VertexContext
}

func NewContext(vertexCtx *types2.VertexContext) *Context {
	return &Context{
		vertexCtx: vertexCtx,
	}
}

func (ctx *Context) AddListener(listener types2.Listener) {
	ctx.vertexCtx.AddListener(listener)
}

func (ctx *Context) RemoveListener(listener types2.Listener) {
	ctx.vertexCtx.RemoveListener(listener)
}

func (ctx *Context) DispatchEvent(event interface{}) {
	ctx.vertexCtx.DispatchEvent(event)
}
