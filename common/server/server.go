package server

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/common"
	"github.com/vertex-center/vertex/common/event"
	"github.com/vertex-center/vertex/common/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
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

	return &Server{
		id:  id,
		url: u,
		ctx: ctx,
		Router: router.New(info,
			router.WithMiddleware(cors.New(cfg)),
			router.WithMiddleware(requestID()),
			router.WithMiddleware(correlationID()),
			router.WithMiddleware(logger(u)),
			router.WithMiddleware(gin.Recovery()),
		),
	}
}

func (s *Server) StartAsync() chan error {
	exitChan := make(chan error)
	go func() {
		defer close(exitChan)
		exitChan <- s.Router.Start(":" + s.url.Port())
	}()

	s.waitServerReady()

	s.ctx.DispatchEvent(event.ServerLoad{})
	s.ctx.DispatchEvent(event.ServerStart{})

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
