package service

import (
	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"

	"github.com/vertex-center/uuid"
)

type proxyService struct {
	proxyAdapter port.ProxyAdapter
}

func NewProxyService(proxyAdapter port.ProxyAdapter) port.ProxyService {
	return &proxyService{
		proxyAdapter: proxyAdapter,
	}
}

func (s *proxyService) GetRedirects() types.ProxyRedirects {
	return s.proxyAdapter.GetRedirects()
}

func (s *proxyService) GetRedirectByHost(host string) *types.ProxyRedirect {
	return s.proxyAdapter.GetRedirectByHost(host)
}

func (s *proxyService) AddRedirect(redirect types.ProxyRedirect) error {
	id := uuid.New()

	if redirect.Source == "" {
		return types.ErrSourceInvalid
	}

	if redirect.Target == "" {
		return types.ErrTargetInvalid
	}

	return s.proxyAdapter.AddRedirect(id, redirect)
}

func (s *proxyService) RemoveRedirect(id uuid.UUID) error {
	return s.proxyAdapter.RemoveRedirect(id)
}
