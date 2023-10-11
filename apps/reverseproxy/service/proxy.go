package service

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/google/uuid"
	"github.com/vertex-center/vertex/apps/reverseproxy/types"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/pkg/router"
	"github.com/vertex-center/vlog"
)

type ProxyService struct {
	proxyAdapter types.ProxyAdapterPort
}

func NewProxyService(proxyAdapter types.ProxyAdapterPort) *ProxyService {
	return &ProxyService{
		proxyAdapter: proxyAdapter,
	}
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

func (s *ProxyService) HandleProxy(c *router.Context) {
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
