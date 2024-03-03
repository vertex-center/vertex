package app

import (
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/pkg/event"
)

type Context struct {
	vertexCtx *common.VertexContext
}

func NewContext(vertexCtx *common.VertexContext) *Context {
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

func (ctx *Context) DispatchEvent(event event.Event) {
	ctx.vertexCtx.DispatchEvent(event)
}

func (ctx *Context) DispatchEventWithErr(event event.Event) error {
	return ctx.vertexCtx.DispatchEventWithErr(event)
}

func (ctx *Context) About() common.About {
	return ctx.vertexCtx.About()
}

func (ctx *Context) Kernel() bool {
	return ctx.vertexCtx.Kernel()
}
