package adapter

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/api"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type sshKernelApiAdapter struct {
	client *api.KernelClient
}

func NewSshKernelApiAdapter() port.SshAdapter {
	return &sshKernelApiAdapter{
		client: api.NewAdminKernelClient(),
	}
}

func (a *sshKernelApiAdapter) GetAll() ([]types.PublicKey, error) {
	return a.client.GetSSHKeys(context.Background())
}

func (a *sshKernelApiAdapter) Add(key string, username string) error {
	return a.client.AddSSHKey(context.Background(), key, username)
}

func (a *sshKernelApiAdapter) Remove(fingerprint string, username string) error {
	return a.client.DeleteSSHKey(context.Background(), fingerprint, username)
}

func (a *sshKernelApiAdapter) GetUsers() ([]user.User, error) {
	return a.client.GetSSHUsers(context.Background())
}
