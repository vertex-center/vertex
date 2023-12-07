package server

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
	"github.com/wI2L/fizz/openapi"
)

var InternetOK = false

type Server struct {
	id     string
	url    *url.URL
	ctx    *types.VertexContext
	Router *router.Router
}

func New(id string, info *openapi.Info, u *url.URL, ctx *types.VertexContext) *Server {
	gin.SetMode(gin.ReleaseMode)

	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("Authorization")

	return &Server{
		id:  id,
		url: u,
		ctx: ctx,
		Router: router.New(info,
			router.WithMiddleware(cors.New(cfg)),
			router.WithMiddleware(ginutils.ErrorHandler()),
			router.WithMiddleware(ginutils.Logger(id, u.String())),
			router.WithMiddleware(gin.Recovery()),
		),
	}
}

func (s *Server) StartAsync() chan error {
	exitChan := make(chan error)
	go func() {
		defer close(exitChan)
		log.Info("starting server", vlog.String("port", s.url.Port()))
		exitChan <- s.Router.Start(":" + s.url.Port())
	}()

	s.waitInternet()
	s.waitServerReady()

	s.ctx.DispatchEvent(types.EventServerLoad{})
	s.ctx.DispatchEvent(types.EventServerStart{})

	return exitChan
}

func (s *Server) Stop() {
	s.ctx.DispatchEvent(types.EventServerStop{})

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := s.Router.Stop(ctx)
	if err != nil {
		log.Error(err)
	}
}

func (s *Server) waitInternet() {
	if InternetOK {
		return
	}

	log.Info("waiting for internet connection...")

	timeout, cancelTimeout := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelTimeout()
	err := net.WaitInternetConn(timeout)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	InternetOK = true
	log.Info("internet connection established")
}

func (s *Server) waitServerReady() {
	pingURL := fmt.Sprintf("%s/ping", s.url.String())

	log.Info("waiting for router to be ready...", vlog.String("ping_url", pingURL))

	timeout, cancelTimeout := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelTimeout()
	err := net.Wait(timeout, pingURL)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("router is ready")
}
