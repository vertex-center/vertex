package router

import "github.com/gin-gonic/gin"

type Group struct {
	*gin.RouterGroup
}

func (g *Group) Group(path string, handlers ...HandlerFunc) *Group {
	return &Group{
		RouterGroup: g.RouterGroup.Group(path, wrapHandlers(handlers...)...),
	}
}

func (g *Group) GET(path string, handlers ...HandlerFunc) {
	g.RouterGroup.GET(path, wrapHandlers(handlers...)...)
}

func (g *Group) POST(path string, handlers ...HandlerFunc) {
	g.RouterGroup.POST(path, wrapHandlers(handlers...)...)
}

func (g *Group) PUT(path string, handlers ...HandlerFunc) {
	g.RouterGroup.PUT(path, wrapHandlers(handlers...)...)
}

func (g *Group) PATCH(path string, handlers ...HandlerFunc) {
	g.RouterGroup.PATCH(path, wrapHandlers(handlers...)...)
}

func (g *Group) DELETE(path string, handlers ...HandlerFunc) {
	g.RouterGroup.DELETE(path, wrapHandlers(handlers...)...)
}

func (g *Group) OPTIONS(path string, handlers ...HandlerFunc) {
	g.RouterGroup.OPTIONS(path, wrapHandlers(handlers...)...)
}

func (g *Group) HEAD(path string, handlers ...HandlerFunc) {
	g.RouterGroup.HEAD(path, wrapHandlers(handlers...)...)
}

func (g *Group) Any(path string, handlers ...HandlerFunc) {
	g.RouterGroup.Any(path, wrapHandlers(handlers...)...)
}
