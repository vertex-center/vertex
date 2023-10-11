package app

import "github.com/vertex-center/vertex/types"

type Context struct {
	vertexCtx *types.VertexContext
}

func NewContext(vertexCtx *types.VertexContext) *Context {
	return &Context{
		vertexCtx: vertexCtx,
	}
}

func (ctx *Context) AddListener(listener types.Listener) {
	ctx.vertexCtx.AddListener(listener)
}

func (ctx *Context) RemoveListener(listener types.Listener) {
	ctx.vertexCtx.RemoveListener(listener)
}

func (ctx *Context) DispatchEvent(event interface{}) {
	ctx.vertexCtx.DispatchEvent(event)
}
