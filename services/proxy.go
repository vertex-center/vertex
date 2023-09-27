package services

import (
	"github.com/google/uuid"
	"github.com/vertex-center/vertex/types"
)

type ProxyService struct {
	proxyAdapter types.ProxyAdapterPort
}

func NewProxyService(proxyAdapter types.ProxyAdapterPort) ProxyService {
	s := ProxyService{
		proxyAdapter: proxyAdapter,
	}
	return s
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
