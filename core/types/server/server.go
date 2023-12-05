package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/net"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

var InternetOK = false

type Server struct {
	id     string
	addr   string
	ctx    *types.VertexContext
	Router *router.Router
}

func New(id, addr string, ctx *types.VertexContext) *Server {
	gin.SetMode(gin.ReleaseMode)

	if addr == "" || addr == ":" {
		log.Error(errors.New("server address is empty"), vlog.String("id", id))
		os.Exit(1)
	}

	s := Server{
		id:     id,
		addr:   addr,
		ctx:    ctx,
		Router: router.New(),
	}
	s.initRouter()
	return &s
}

func (s *Server) initRouter() {
	cfg := cors.DefaultConfig()
	cfg.AllowAllOrigins = true
	cfg.AddAllowHeaders("Authorization")

	s.Router.Use(cors.New(cfg))
	s.Router.Use(ginutils.ErrorHandler())
	s.Router.Use(ginutils.Logger(s.id, s.addr))
	s.Router.Use(gin.Recovery())
}

func (s *Server) StartAsync() chan error {
	exitChan := make(chan error)
	go func() {
		defer close(exitChan)
		log.Info("starting server", vlog.String("addr", s.addr))
		exitChan <- s.Router.Start(s.addr)
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
	log.Info("waiting for router to be ready...")

	timeout, cancelTimeout := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelTimeout()
	url := fmt.Sprintf("http://%s%s/api/ping", config.Current.Host, s.addr)
	err := net.Wait(timeout, url)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	log.Info("router is ready")
}
