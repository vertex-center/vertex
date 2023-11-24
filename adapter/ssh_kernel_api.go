package adapter

import (
	"context"

	"github.com/carlmjohnson/requests"
	"github.com/vertex-center/vertex/config"
	"github.com/vertex-center/vertex/core/port"
	"github.com/vertex-center/vertex/core/service"
	"github.com/vertex-center/vertex/core/types"
	"github.com/vertex-center/vertex/core/types/api"
	"github.com/vertex-center/vertex/handler"
	"github.com/vertex-center/vertex/pkg/user"
)

type SshKernelApiAdapter struct {
	config requests.Config
}

func NewSshKernelApiAdapter() port.SshAdapter {
	return &SshKernelApiAdapter{
		config: func(rb *requests.Builder) {
			rb.BaseURL(config.Current.KernelURL())
		},
	}
}

func (a *SshKernelApiAdapter) GetAll() ([]types.PublicKey, error) {
	var keys []types.PublicKey
	var apiError api.Error

	err := requests.New(a.config).
		Path("/api/security/ssh").
		ToJSON(&keys).
		ErrorJSON(&apiError).
		Fetch(context.Background())
	return keys, err
}

func (a *SshKernelApiAdapter) Add(key string, username string) error {
	var apiError api.Error

	err := requests.New(a.config).
		Path("/api/security/ssh").
		BodyJSON(&handler.AddSSHKeyBody{
			AuthorizedKey: key,
			Username:      username,
		}).
		Post().
		ErrorJSON(&apiError).
		Fetch(context.Background())

	if apiError.Code == api.ErrUserNotFound {
		return service.ErrUserNotFound
	}
	return err
}

func (a *SshKernelApiAdapter) Remove(fingerprint string, username string) error {
	var apiError api.Error

	err := requests.New(a.config).
		Pathf("/api/security/ssh").
		BodyJSON(&handler.DeleteSSHKeyBody{
			Fingerprint: fingerprint,
			Username:    username,
		}).
		Delete().
		ErrorJSON(&apiError).
		Fetch(context.Background())

	if apiError.Code == api.ErrUserNotFound {
		return service.ErrUserNotFound
	}
	return err
}

func (a *SshKernelApiAdapter) GetUsers() ([]user.User, error) {
	var users []user.User
	var apiError api.Error

	err := requests.New(a.config).
		Path("/api/security/ssh/users").
		ToJSON(&users).
		ErrorJSON(&apiError).
		Fetch(context.Background())
	return users, err
}
