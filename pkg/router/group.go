package router

import (
	"github.com/vertex-center/vertex/pkg/router/oapi"
	"github.com/wI2L/fizz"
)

type Group struct {
	*fizz.RouterGroup
}

func (g *Group) Group(path, name, description string, handlers ...HandlerFunc) *Group {
	return &Group{
		RouterGroup: g.RouterGroup.Group(path, name, description, wrapHandlers(handlers...)...),
	}
}

func (g *Group) GET(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.GET(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) POST(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.POST(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) PUT(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.PUT(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) PATCH(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.PATCH(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) DELETE(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.DELETE(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) OPTIONS(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.OPTIONS(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) HEAD(path string, infos []oapi.Info, handlers ...HandlerFunc) {
	g.RouterGroup.HEAD(path, oapi.WrapInfos(infos...), wrapHandlers(handlers...)...)
}

func (g *Group) Use(handlers ...HandlerFunc) {
	g.RouterGroup.Use(wrapHandlers(handlers...)...)
}
