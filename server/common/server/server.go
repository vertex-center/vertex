package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/loopfz/gadgeto/tonic"
	"github.com/vertex-center/vertex/server/common"
	"github.com/vertex-center/vertex/server/common/event"
	"github.com/vertex-center/vertex/server/common/log"
	"github.com/vertex-center/vertex/server/pkg/net"
	"github.com/vertex-center/vertex/server/pkg/router"
	"github.com/vertex-center/vlog"
	"github.com/wI2L/fizz"
	"github.com/wI2L/fizz/openapi"
)

type Server struct {
	id     string
	url    *url.URL
	ctx    *common.VertexContext
	Router *router.Router
}

func New(id string, info *openapi.Info, u *url.URL, ctx *common.VertexContext) *Server {
	gin.SetMode(gin.ReleaseMode)

	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("Authorization")
	cfg.AddAllowHeaders("X-Request-ID")
	cfg.AddAllowHeaders("X-Correlation-ID")

	r := router.New(info,
		router.WithMiddleware(cors.New(cfg)),
		router.WithMiddleware(requestID()),
		router.WithMiddleware(correlationID()),
		router.WithMiddleware(logger(u, id)),
		router.WithMiddleware(gin.Recovery()),
	)

	r.GET("/api/about", []fizz.OperationOption{
		fizz.ID("getAbout"),
		fizz.Summary("Get server info"),
	}, tonic.Handler(func(c *gin.Context) (*common.About, error) {
		a := ctx.About()
		return &a, nil
	}, http.StatusOK))

	return &Server{
		id:     id,
		url:    u,
		ctx:    ctx,
		Router: r,
	}
}

func (s *Server) StartAsync() chan error {
	exitChan := make(chan error)
	go func() {
		defer close(exitChan)
		log.Info("server starting", vlog.String("url", s.url.String()), vlog.String("port", s.url.Port()))
		exitChan <- s.Router.Start("0.0.0.0:" + s.url.Port())
	}()

	s.waitServerReady()

	s.ctx.DispatchEvent(event.ServerLoad{})
	s.ctx.DispatchEvent(event.ServerStart{})

	log.Info("server started", vlog.String("port", s.url.Port()))

	return exitChan
}

func (s *Server) Stop() {
	s.ctx.DispatchEvent(event.ServerStop{})

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := s.Router.Stop(ctx)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) waitServerReady() {
	pingURL := fmt.Sprintf("%s/ping", s.url.String())
	timeout, cancelTimeout := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelTimeout()
	err := net.Wait(timeout, pingURL)
	if err != nil {
		panic(err)
	}
}
