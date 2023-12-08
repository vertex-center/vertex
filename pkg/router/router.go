package router

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
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
