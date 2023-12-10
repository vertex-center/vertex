package adapter

import (
	"context"

	"github.com/vertex-center/vertex/apps/admin/api"
	"github.com/vertex-center/vertex/apps/admin/core/port"
	"github.com/vertex-center/vertex/apps/admin/core/types"
	"github.com/vertex-center/vertex/pkg/user"
)

type sshKernelApiAdapter struct{}

func NewSshKernelApiAdapter() port.SshAdapter {
	return &sshKernelApiAdapter{}
}

func (a *sshKernelApiAdapter) GetAll(ctx context.Context) ([]types.PublicKey, error) {
	cli := api.NewAdminKernelClient(ctx)
	return cli.GetSSHKeys(ctx)
}

func (a *sshKernelApiAdapter) Add(ctx context.Context, key string, username string) error {
	cli := api.NewAdminKernelClient(ctx)
	return cli.AddSSHKey(ctx, key, username)
}

func (a *sshKernelApiAdapter) Remove(ctx context.Context, fingerprint string, username string) error {
	cli := api.NewAdminKernelClient(ctx)
	return cli.DeleteSSHKey(ctx, fingerprint, username)
}

func (a *sshKernelApiAdapter) GetUsers(ctx context.Context) ([]user.User, error) {
	cli := api.NewAdminKernelClient(ctx)
	return cli.GetSSHUsers(ctx)
}
