package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/vertex-center/vertex/pkg/log"
	"github.com/vertex-center/vertex/types"
	"github.com/vertex-center/vlog"
)

type ProxyKernelService struct {
	proxies types.ProxyRedirects
}

func NewProxyKernelService() ProxyKernelService {
	return ProxyKernelService{
		proxies: types.ProxyRedirects{},
	}
}

func (s *ProxyKernelService) SetRedirects(proxies types.ProxyRedirects) {
	s.proxies = proxies
}

func (s *ProxyKernelService) HandleProxy(c *gin.Context) {
	host := c.Request.Host

	redirect := s.getRedirectByHost(host)
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

func (s *ProxyKernelService) getRedirectByHost(host string) *types.ProxyRedirect {
	for _, redirect := range s.proxies {
		if redirect.Source == host {
			return &redirect
		}
	}
	return nil
}
