package adapter

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/api"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type SshKernelApiAdapter struct {
	client *api.KernelClient
}

func NewSshKernelApiAdapter() port.SshAdapter {
	return &SshKernelApiAdapter{
		client: api.NewAdminKernelClient(),
	}
}

func (a *SshKernelApiAdapter) GetAll() ([]types.PublicKey, error) {
	return a.client.GetSSHKeys(context.Background())
}

func (a *SshKernelApiAdapter) Add(key string, username string) error {
	return a.client.AddSSHKey(context.Background(), key, username)
}

func (a *SshKernelApiAdapter) Remove(fingerprint string, username string) error {
	return a.client.DeleteSSHKey(context.Background(), fingerprint, username)
}

func (a *SshKernelApiAdapter) GetUsers() ([]user.User, error) {
	return a.client.GetSSHUsers(context.Background())
}
