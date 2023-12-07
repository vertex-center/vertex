package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router/oapi"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Router struct {
	*fizz.Fizz

	server *http.Server
}

type Option func(*gin.Engine)

func WithMiddleware(middleware ...gin.HandlerFunc) Option {
	return func(r *gin.Engine) {
		r.Use(middleware...)
	}
}

func New(info *openapi.Info, opts ...Option) *Router {
	e := gin.New()
	for _, opt := range opts {
		opt(e)
	}
	f := fizz.NewFromEngine(e)
	if info != nil {
		f.GET("/openapi.yaml", nil, f.OpenAPI(info, "yaml"))
		f.GET("/openapi.json", nil, f.OpenAPI(info, "json"))
	}
	if len(f.Errors()) > 0 {
		log.Error(errors.Join(f.Errors()...))
	}
	return &Router{
		Fizz: f,
	}
}

func (r *Router) Start(addr string) error {
	r.server = &http.Server{
		Addr:    addr,
		Handler: r.Fizz,
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

func (r *Router) Group(path, name, description string, handlers ...gin.HandlerFunc) *Group {
	return &Group{
		RouterGroup: r.Fizz.Group(path, name, description, handlers...),
	}
}

func (r *Router) GET(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.GET(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) POST(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.POST(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) PUT(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.PUT(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) PATCH(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.PATCH(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) DELETE(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.DELETE(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) OPTIONS(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.OPTIONS(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) HEAD(path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.HEAD(path, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) Handle(method, path string, infos []oapi.Info, handlers ...gin.HandlerFunc) {
	r.RouterGroup.Handle(path, method, oapi.WrapInfos(infos...), handlers...)
}

func (r *Router) Any(path string, handlers ...gin.HandlerFunc) {
	r.Engine().Any(path, handlers...)
}
