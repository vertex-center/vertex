package service

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type SshService struct {
	adapter port.SshAdapter
}

func NewSshService(sshAdapter port.SshAdapter) port.SshService {
	return &SshService{
		adapter: sshAdapter,
	}
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

func (s *SshService) GetUsers() ([]user.User, error) {
	return s.adapter.GetUsers()
}
