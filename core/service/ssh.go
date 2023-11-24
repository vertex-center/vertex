package service

import (
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/types"
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

func (s *SshService) Add(key string, username string) error {
	return s.adapter.Add(key, username)
}

func (s *SshService) Delete(fingerprint string, username string) error {
	return s.adapter.Remove(fingerprint, username)
}

func (s *SshService) GetUsers() ([]string, error) {
	users, err := s.adapter.GetUsers()
	if err != nil {
		return nil, err
	}

	var res []string
	for _, u := range users {
		res = append(res, u.Name)
	}

	return res, nil
}
