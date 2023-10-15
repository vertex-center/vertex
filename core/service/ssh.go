package service

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
)

type SshService struct {
	adapter port.SshAdapter
}

func NewSshService(sshAdapter port.SshAdapter) SshService {
	s := SshService{
		adapter: sshAdapter,
	}
	return s
}

func (s *SshService) GetAll() ([]types.PublicKey, error) {
	return s.adapter.GetAll()
}

func (s *SshService) Add(key string) error {
	return s.adapter.Add(key)
}

func (s *SshService) Delete(fingerprint string) error {
	return s.adapter.Remove(fingerprint)
}
