package service

import (
	"context"

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

func (s *sshService) GetAll(ctx context.Context) ([]types.PublicKey, error) {
	return s.adapter.GetAll(ctx)
}

func (s *sshService) Add(ctx context.Context, key string, username string) error {
	return s.adapter.Add(ctx, key, username)
}

func (s *sshService) Delete(ctx context.Context, fingerprint string, username string) error {
	return s.adapter.Remove(ctx, fingerprint, username)
}

func (s *sshService) GetUsers(ctx context.Context) ([]string, error) {
	users, err := s.adapter.GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, u := range users {
		res = append(res, u.Name)
	}

	return res, nil
}
