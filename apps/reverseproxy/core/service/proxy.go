package service

import (
	"fmt"

	"github.com/vertex-center/vertex/apps/reverseproxy/core/port"
	"github.com/vertex-center/vertex/apps/reverseproxy/core/types"

	"github.com/google/uuid"
)

type ProxyService struct {
	proxyAdapter port.ProxyAdapter
}

func NewProxyService(proxyAdapter port.ProxyAdapter) port.ProxyService {
	return &ProxyService{
		proxyAdapter: proxyAdapter,
	}
}

func (s *ProxyService) GetRedirects() types.ProxyRedirects {
	return s.proxyAdapter.GetRedirects()
}

func (s *ProxyService) GetRedirectByHost(host string) *types.ProxyRedirect {
	return s.proxyAdapter.GetRedirectByHost(host)
}

func (s *ProxyService) AddRedirect(redirect types.ProxyRedirect) error {
	id := uuid.New()

	if redirect.Source == "" {
		return fmt.Errorf("Failed to add redirect as source is empty.")
	}

	if redirect.Target == "" {
		return fmt.Errorf("Failed to add redirect as target is empty.")
	}

	return s.proxyAdapter.AddRedirect(id, redirect)
}

func (s *ProxyService) RemoveRedirect(id uuid.UUID) error {
	return s.proxyAdapter.RemoveRedirect(id)
}
