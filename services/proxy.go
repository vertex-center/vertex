package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/pkg/ginutils"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

var (
	ErrProxyAlreadyRunning = errors.New("a proxy is already running, cannot start a new one")
)

type ProxyService struct {
	server       *http.Server
	proxyAdapter types.ProxyAdapterPort
}

func NewProxyService(proxyAdapter types.ProxyAdapterPort) ProxyService {
	s := ProxyService{
		proxyAdapter: proxyAdapter,
	}
	return s
}

func (s *ProxyService) Start() error {
	if s.server != nil {
		return ErrProxyAlreadyRunning
	}

	r := gin.New()
	r.Use(cors.Default())
	r.Use(ginutils.Logger("PROX"))
	r.Use(gin.Recovery())
	r.Any("/*path", s.handleProxy)

	s.server = &http.Server{
		Addr:    ":80",
		Handler: r,
	}

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error(err)
			return
		}
	}()

	return nil
}

func (s *ProxyService) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	s.server = nil
	return nil
}

func (s *ProxyService) GetRedirects() types.ProxyRedirects {
	return s.proxyAdapter.GetRedirects()
}

func (s *ProxyService) AddRedirect(redirect types.ProxyRedirect) error {
	id := uuid.New()
	return s.proxyAdapter.AddRedirect(id, redirect)
}

func (s *ProxyService) RemoveRedirect(id uuid.UUID) error {
	return s.proxyAdapter.RemoveRedirect(id)
}

func (s *ProxyService) handleProxy(c *gin.Context) {
	host := c.Request.Host

	redirect := s.proxyAdapter.GetRedirectByHost(host)
	if redirect == nil {
		log.Warn("this host is not registered in the reverse proxy",
			vlog.String("host", host),
		)
		return
	}

	target, err := url.Parse(redirect.Target)
	if err != nil {
		log.Error(err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, request *http.Request, err error) {
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Error(err)
		}
	}
	proxy.Director = func(request *http.Request) {
		request.Header = c.Request.Header
		request.Host = target.Host
		request.URL.Scheme = target.Scheme
		request.URL.Host = target.Host
		request.URL.Path = c.Param("path")
	}
	proxy.ServeHTTP(c.Writer, c.Request)
}
