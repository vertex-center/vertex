package service

import (
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
)

type sshService struct {
	adapter port.SshAdapter
}

func NewSshService(sshAdapter port.SshAdapter) port.SshService {
	return &sshService{
		adapter: sshAdapter,
	}
}

func (s *sshService) GetAll() ([]types.PublicKey, error) {
	return s.adapter.GetAll()
}

func (s *sshService) Add(key string, username string) error {
	return s.adapter.Add(key, username)
}

func (s *sshService) Delete(fingerprint string, username string) error {
	return s.adapter.Remove(fingerprint, username)
}

func (s *sshService) GetUsers() ([]string, error) {
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
