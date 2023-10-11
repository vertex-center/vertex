package app

import (
	"github.com/gin-contrib/sse"
	"github.com/vertex-center/vertex/pkg/router"
)

type App interface {
	Initialize(registry *AppsRegistry) error
	Uninitialize(registry *AppsRegistry) error

	Name() string

	OnEvent(e interface{})
}

type Router interface {
	AddRoutes(r *router.Group)
}

type Service interface {
	OnEvent(e interface{})
}

func HeadersSSE(c *router.Context) {
	c.Writer.Header().Set("Content-Type", sse.ContentType)
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}
