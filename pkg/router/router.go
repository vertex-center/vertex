package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Router struct {
	*gin.Engine
	server *http.Server
}

func New() *Router {
	return &Router{
		Engine: gin.New(),
	}
}

func (r *Router) Start(addr string) error {
	r.server = &http.Server{
		Addr:    addr,
		Handler: r.Engine,
	}
	err := r.server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	return err
}

// Stop gracefully shuts down the server. It will throw ErrFailedToStopServer
// if the server fails to stop.
func (r *Router) Stop(ctx context.Context) error {
	if r.server == nil {
		return nil
	}
	err := r.server.Shutdown(ctx)
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	} else if err != nil {
		err = fmt.Errorf("%w: %w", ErrFailedToStopServer, err)
	}
	r.server = nil
	return err
}

func (r *Router) Group(path string, handlers ...HandlerFunc) *Group {
	return &Group{
		RouterGroup: r.Engine.Group(path, wrapHandlers(handlers...)...),
	}
}

func (r *Router) GET(path string, handlers ...HandlerFunc) {
	r.RouterGroup.GET(path, wrapHandlers(handlers...)...)
}

func (r *Router) POST(path string, handlers ...HandlerFunc) {
	r.RouterGroup.POST(path, wrapHandlers(handlers...)...)
}

func (r *Router) PUT(path string, handlers ...HandlerFunc) {
	r.RouterGroup.PUT(path, wrapHandlers(handlers...)...)
}

func (r *Router) PATCH(path string, handlers ...HandlerFunc) {
	r.RouterGroup.PATCH(path, wrapHandlers(handlers...)...)
}

func (r *Router) DELETE(path string, handlers ...HandlerFunc) {
	r.RouterGroup.DELETE(path, wrapHandlers(handlers...)...)
}

func (r *Router) OPTIONS(path string, handlers ...HandlerFunc) {
	r.RouterGroup.OPTIONS(path, wrapHandlers(handlers...)...)
}

func (r *Router) HEAD(path string, handlers ...HandlerFunc) {
	r.RouterGroup.HEAD(path, wrapHandlers(handlers...)...)
}

func (r *Router) Handle(method, path string, handlers ...HandlerFunc) {
	r.RouterGroup.Handle(method, path, wrapHandlers(handlers...)...)
}

func (r *Router) Any(path string, handlers ...HandlerFunc) {
	r.RouterGroup.Any(path, wrapHandlers(handlers...)...)
}
