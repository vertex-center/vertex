package services

import (
	"github.com/vertex-center/vertex/types"
)

type SshService struct {
	adapter types.SshAdapterPort
}

func NewSSHService(sshAdapter types.SshAdapterPort) SshService {
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
